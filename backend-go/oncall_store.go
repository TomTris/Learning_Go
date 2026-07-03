package main

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OnCallStore interface {
	Create(ctx context.Context, entry OnCallShiftEntry) (OnCallShiftEntry, error)
	CurrentOnCall(ctx context.Context, service string) (string, error)
	CurrentOnCallAll(ctx context.Context) ([]OnCallShiftEntry, error)
	ListOnCalls(ctx context.Context, from *time.Time, to *time.Time) ([]OnCallShiftEntry, error) //minute-exact
	UpdateOnCall(ctx context.Context, updatedEntry OnCallShiftEntry) (OnCallShiftEntry, error)
}

func NewOnCallStore(ctx context.Context, db *mongo.Database) (OnCallStore, error) {
	if db == nil {
		slog.Info("use in-memory store for UserStore")
		return NewInMemoryOnCallStore()
	}

	slog.Info("use MongoStore for UserStore")
	return NewMongoOnCallStore(ctx, db.Collection(CollectionOnCallShifts))
}
