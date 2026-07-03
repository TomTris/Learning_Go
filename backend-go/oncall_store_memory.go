package main

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type InMemoryOnCallStore struct {
	mu            sync.RWMutex
	OnCallEntries map[string]OnCallShiftEntry
	currentID     int
}

func NewInMemoryOnCallStore() (*InMemoryOnCallStore, error) {
	s := InMemoryOnCallStore{
		OnCallEntries: make(map[string]OnCallShiftEntry), // id:Entry
		currentID:     0,
	}
	return &s, nil
}

func (s *InMemoryOnCallStore) Create(ctx context.Context, entry OnCallShiftEntry) (OnCallShiftEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentID++
	entry.ID = OnCallShiftEntryIDPrefix + strconv.Itoa(s.currentID)
	s.OnCallEntries[entry.ID] = entry
	return entry, nil
}

func (s *InMemoryOnCallStore) CurrentOnCall(ctx context.Context, service string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	for _, each := range s.OnCallEntries {
		// startsat <= now  AND  now > EndsAt
		if each.Service == service && !each.StartsAt.After(now) && each.EndsAt.After(now) {
			return each.Username, nil
		}
	}
	return "", OnCallShiftEntryNotFound
}

func (s *InMemoryOnCallStore) CurrentOnCallAll(ctx context.Context) ([]OnCallShiftEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	entries := []OnCallShiftEntry{}
	for _, e := range s.OnCallEntries {
		if !e.StartsAt.After(now) && e.EndsAt.After(now) { // start <= now < end
			entries = append(entries, e)
		}
	}
	return entries, nil
}

func (s *InMemoryOnCallStore) ListOnCalls(ctx context.Context, from *time.Time, to *time.Time) ([]OnCallShiftEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var lo, hi *time.Time
	if from != nil {
		f := from.Truncate(time.Minute) // start of the minute
		lo = &f
	}
	if to != nil {
		t := to.Truncate(time.Minute).Add(time.Minute - time.Nanosecond) // end of the minute
		hi = &t
	}

	entries := []OnCallShiftEntry{}
	for _, each := range s.OnCallEntries {
		// overlap requires: end >= from AND start <= to
		if lo != nil && each.EndsAt.Before(*lo) {
			continue // end < from -> no overlap
		}
		if hi != nil && each.StartsAt.After(*hi) {
			continue // start > to -> no overlap
		}
		entries = append(entries, each)
	}
	return entries, nil
}

func (s *InMemoryOnCallStore) UpdateOnCall(ctx context.Context, updatedEntry OnCallShiftEntry) (OnCallShiftEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.OnCallEntries[updatedEntry.ID]; !ok {
		return OnCallShiftEntry{}, OnCallShiftEntryNotFound
	}

	s.OnCallEntries[updatedEntry.ID] = updatedEntry
	return updatedEntry, nil
}
