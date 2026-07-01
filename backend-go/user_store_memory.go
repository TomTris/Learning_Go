package main

import (
	"context"
	"strconv"
	"sync"
)

type MemoryUserStore struct {
	mu        sync.RWMutex
	users     map[string]User // username - User
	currentID int
}

func NewMemoryUserStore() (*MemoryUserStore, error) {
	m := make(map[string]User)
	return &MemoryUserStore{users: m}, nil
}

func NewMemoryUserStoreWithSeed(seed []User) *MemoryUserStore {
	m := make(map[string]User, len(seed))
	for _, u := range seed {
		m[u.Username] = u
	}
	return &MemoryUserStore{users: m}
}

func (s *MemoryUserStore) Create(_ context.Context, u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.users[u.Username]
	if ok == true {
		return User{}, ErrUserAlreadyExist
	}

	s.currentID++
	ID := UserIDPrefix + strconv.Itoa(s.currentID)
	u.ID = ID
	s.users[u.Username] = u
	return u, nil
}

func (s *MemoryUserStore) GetByUsername(_ context.Context, username string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[username]
	if ok == false {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
