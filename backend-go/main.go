package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database

func getMongoDatabase(conf Config) *mongo.Database {
	if conf.ConnectionString == "" {
		slog.Info("HANDOFF_CONNECT_STRING is empty, use Memory store only")
		return nil
	}

	slog.Info("using mongo store", "db", conf.DatabaseName)
	client, err := mongo.Connect(options.Client().ApplyURI(conf.ConnectionString))
	if err != nil {
		log.Fatal("can't connect to db via HANDOFF_CONNECT_STRING")
	}

	db := client.Database(conf.DatabaseName)
	return db
}

func mongoNextID(ctx context.Context, nameInCollectionCounter string, prefix string) (string, error) {
	if db == nil {
		slog.Error("db is nil")
		return "", errors.New("db is nil")
	}
	col := db.Collection(CollectionCounters)
	filter := bson.M{"_id": nameInCollectionCounter}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return "", err
	}
	return prefix + strconv.Itoa(result.Seq), nil
}

func main() {
	// init metrics
	promRegistry := prometheus.NewRegistry()
	httpMetrics := NewHttpMetrics(promRegistry)
	metricRegistry := NewMetricRegistry(promRegistry)
	incidentStoreMetric := NewIncidentStoreMetric(promRegistry)

	// init Registry (Websocket connection)
	registry := NewRegistry(metricRegistry)
	go registry.run()
	defer close(registry.done)

	// init flagHandler
	flagHandler := FlagHandler{store: CreateFlagStore()}

	// init config and mongoClient
	config := loadConfig()
	db = getMongoDatabase(config)

	// init onCallHandler
	onCallStore, err := NewOnCallStore(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
	onCallHandler := &OnCallHandler{Store: onCallStore}
	// Init IncidentHandler and its store
	incidentStore, err := NewIncidentStore(db)
	if err != nil {
		log.Fatalf("%v", err) // for an error
	}
	instrumentedIncidentStore := InstrumentedIncidentStore{
		inner:   incidentStore,
		metrics: incidentStoreMetric,
	}
	incHandler := IncidentHandler{
		IncidentStore: &instrumentedIncidentStore,
		Registry:      registry,
		FlagEvaluator: &flagHandler.store,
		CurrentOnCall: onCallHandler.Store,
	}

	// init authHandler and its store
	userStore, err := NewUserStore(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
	authHandler := NewAuthHandler(userStore, []byte(config.JWT_SECRET), time.Duration(15*time.Minute))

	// Set router
	router := getRouter(&incHandler, &flagHandler, authHandler, onCallHandler, db, promRegistry, httpMetrics)

	// run server
	srv := http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}
	go func() {
		slog.Info(fmt.Sprintf("server starting http://127.0.0.1%s/", srv.Addr))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// greaceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Info("server shut down in <= 10 sec")
	srv.Shutdown(ctx)
	slog.Info("server shut down gracefully")
}
