package main

import (
	"context"
	"testing"
)

// userStoreContract exercises the UserStore interface. Any implementation
// must pass it. newStore must return a fresh, empty store on each call —
// including a reset ID counter, so the ID-sequence assertion holds.
func userStoreContract(t *testing.T, newStore func() UserStore) {
	ctx := context.Background()

	t.Run("Create returns sequential IDs and preserves fields", func(t *testing.T) {
		store := newStore()
		u1, err := store.Create(ctx, User{Username: "tom"})
		if err != nil {
			t.Fatalf("create 1: %v", err)
		}
		u2, err := store.Create(ctx, User{Username: "jerry"})
		if err != nil {
			t.Fatalf("create 2: %v", err)
		}
		if u1.ID != UserIDPrefix+"1" || u2.ID != UserIDPrefix+"2" {
			t.Fatalf("expected sequential IDs, got %q %q", u1.ID, u2.ID)
		}
		if u1.Username != "tom" {
			t.Fatalf("username not preserved, got %q", u1.Username)
		}
	})

	t.Run("GetByUsername returns created user", func(t *testing.T) {
		store := newStore()
		created, err := store.Create(ctx, User{Username: "tom"})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		got, err := store.GetByUsername(ctx, "tom")
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != created.ID || got.Username != "tom" {
			t.Fatalf("mismatch: created %+v, got %+v", created, got)
		}
	})

	t.Run("GetByUsername not found for unknown username", func(t *testing.T) {
		store := newStore()
		_, err := store.GetByUsername(ctx, "nobody")
		if err != ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("Create rejects duplicate username", func(t *testing.T) {
		store := newStore()
		if _, err := store.Create(ctx, User{Username: "tom"}); err != nil {
			t.Fatalf("first create: %v", err)
		}
		_, err := store.Create(ctx, User{Username: "tom"})
		if err != ErrUserAlreadyExist {
			t.Fatalf("expected ErrUserAlreadyExist, got %v", err)
		}
	})

	t.Run("Create allows distinct usernames", func(t *testing.T) {
		store := newStore()
		if _, err := store.Create(ctx, User{Username: "tom"}); err != nil {
			t.Fatalf("create tom: %v", err)
		}
		if _, err := store.Create(ctx, User{Username: "jerry"}); err != nil {
			t.Fatalf("create jerry: %v", err)
		}
	})
}

func TestUserStoreMemoryContract(t *testing.T) {
	userStoreContract(t, func() UserStore {
		NewMemoryUserStore, _ := NewMemoryUserStore()
		return NewMemoryUserStore
	})
}

func setupUserStoreMongoContractEnv(t *testing.T) *MongoUserStore {
	t.Helper()
	config := loadConfig()
	db = getMongoDatabase(config)
	db.Drop(t.Context()) // full drop resets ID counter and removes indexes
	store, err := NewMongoUserStore(t.Context(), db.Collection(CollectionUsers))
	if err != nil {
		t.Fatalf("new mongo user store: %v", err)
	}
	return store
}

func TestUserStoreMongoContract(t *testing.T) {
	userStoreContract(t, func() UserStore {
		return setupUserStoreMongoContractEnv(t)
	})
}
