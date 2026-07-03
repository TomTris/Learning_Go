package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoOnCallStore struct {
	col *mongo.Collection
}

func NewMongoOnCallStore(ctx context.Context, col *mongo.Collection) (*MongoOnCallStore, error) {
	return &MongoOnCallStore{col: col}, nil
}

func (s *MongoOnCallStore) Create(ctx context.Context, entry OnCallShiftEntry) (OnCallShiftEntry, error) {
	id, err := mongoNextID(ctx, CollectionOnCallShifts, OnCallShiftEntryIDPrefix)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("can not get next on-call id: %w", err)
	}
	entry.ID = id

	_, err = s.col.InsertOne(ctx, entry)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("insert on-call error: %w", err)
	}
	return entry, nil
}

func (s *MongoOnCallStore) CurrentOnCall(ctx context.Context, service string) (string, error) {
	now := time.Now()

	filter := bson.M{
		"service":   service,
		"starts_at": bson.M{"$lte": now},
		"ends_at":   bson.M{"$gt": now},
	}

	var entry OnCallShiftEntry
	err := s.col.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", OnCallShiftEntryNotFound
		}
		return "", fmt.Errorf("current on-call query error: %w", err)
	}
	return entry.Username, nil
}

func (s *MongoOnCallStore) CurrentOnCallAll(ctx context.Context) ([]OnCallShiftEntry, error) {
	now := time.Now()
	filter := bson.M{
		"starts_at": bson.M{"$lte": now},
		"ends_at":   bson.M{"$gt": now},
	}
	cursor, err := s.col.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "service", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("current on-call all query error: %w", err)
	}
	entries := []OnCallShiftEntry{}
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("current on-call all decode error: %w", err)
	}
	return entries, nil
}

func (s *MongoOnCallStore) ListOnCalls(ctx context.Context, from *time.Time, to *time.Time) ([]OnCallShiftEntry, error) {
	filter := bson.M{}

	if from != nil {
		// normalize: from -> start of its minute
		f := from.Truncate(time.Minute)
		// on-call overlaps only if it ends at or after `from`: end >= f
		filter["ends_at"] = bson.M{"$gte": f}
	}
	if to != nil {
		// normalize: to -> end of its minute (:59.999999999)
		t := to.Truncate(time.Minute).Add(time.Minute - time.Nanosecond)
		// on-call overlaps only if it starts at or before `to`: start <= t
		filter["starts_at"] = bson.M{"$lte": t}
	}

	cursor, err := s.col.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "starts_at", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list on-calls query error: %w", err)
	}

	entries := []OnCallShiftEntry{}
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("list on-calls decode error: %w", err)
	}
	return entries, nil
}

func (s *MongoOnCallStore) UpdateOnCall(ctx context.Context, updatedEntry OnCallShiftEntry) (OnCallShiftEntry, error) {
	filter := bson.M{"_id": updatedEntry.ID}

	res, err := s.col.ReplaceOne(ctx, filter, updatedEntry)
	if err != nil {
		return OnCallShiftEntry{}, fmt.Errorf("update on-call error: %w", err)
	}
	if res.MatchedCount == 0 {
		return OnCallShiftEntry{}, OnCallShiftEntryNotFound
	}
	return updatedEntry, nil
}
