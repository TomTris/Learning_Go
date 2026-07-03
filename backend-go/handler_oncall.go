package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type CurrentOnCall interface {
	CurrentOnCall(ctx context.Context, service string) (string, error)
}

type OnCallHandler struct {
	Store OnCallStore
}

type OnCallShiftEntry struct {
	ID       string    `json:"id" bson:"_id"`
	Service  string    `json:"service" bson:"service"`
	Username string    `json:"username" bson:"username"`
	StartsAt time.Time `json:"starts_at" bson:"starts_at"`
	EndsAt   time.Time `json:"ends_at" bson:"ends_at"`
}

func (entry *OnCallShiftEntry) Validation() error {
	if strings.TrimSpace(entry.Service) == "" {
		return ErrBadRequest
	}
	if strings.TrimSpace(entry.Username) == "" {
		return ErrBadRequest
	}
	if entry.StartsAt.After(entry.EndsAt) == true {
		return ErrBadRequest
	}
	return nil
}

// TODO: Define error code and have clear return value
func (h *OnCallHandler) CreateShift(r *http.Request) (*AppResponse, *AppError) {
	entry := OnCallShiftEntry{}
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}

	if err := entry.Validation(); err != nil {
		return nil, BadRequest(err)
	}

	created_entry, err := h.Store.Create(r.Context(), entry)
	if err != nil {
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusCreated, created_entry), nil
}

// TODO: Define error code and have clear return value
func (h *OnCallHandler) CurrentOnCall(r *http.Request) (*AppResponse, *AppError) {
	service := r.URL.Query().Get("service")
	if service == "" {
		return nil, BadRequest(ErrServiceRequired)
	}
	username, err := h.Store.CurrentOnCall(r.Context(), service)
	if err != nil {
		return nil, NotFound(err)
	}
	return newAppResponse(http.StatusOK, map[string]string{"username": username}), nil
}

func (h *OnCallHandler) CurrentOnCallAll(r *http.Request) (*AppResponse, *AppError) {
	entries, err := h.Store.CurrentOnCallAll(r.Context())
	if err != nil {
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusOK, entries), nil
}

func (h *OnCallHandler) ListOnCalls(r *http.Request) (*AppResponse, *AppError) {
	startsAt, err := parseTimeQuery(r, "from")
	if err != nil {
		return nil, BadRequest(err)
	}
	endsAt, err := parseTimeQuery(r, "to")
	if err != nil {
		return nil, BadRequest(err)
	}

	// Validate
	if startsAt != nil && endsAt != nil && startsAt.After(*endsAt) {
		return nil, BadRequest(ErrBadRequest)
	}

	entries, storeErr := h.Store.ListOnCalls(r.Context(), startsAt, endsAt)
	if storeErr != nil {
		return nil, InternalServerError(storeErr)
	}
	return newAppResponse(http.StatusOK, entries), nil
}

func (h *OnCallHandler) UpdateShift(r *http.Request) (*AppResponse, *AppError) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		return nil, BadRequest(ErrBadRequest)
	}

	entry := OnCallShiftEntry{}
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}
	entry.ID = id // path is authoritative

	if err := entry.Validation(); err != nil {
		return nil, BadRequest(err)
	}

	updated, err := h.Store.UpdateOnCall(r.Context(), entry)
	if err != nil {
		if errors.Is(err, OnCallShiftEntryNotFound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusOK, updated), nil
}

// parseTimeQuery reads an optional RFC3339 timestamp query param.
// Absent -> (nil, nil). Present but invalid -> (nil, ErrBadRequest).
func parseTimeQuery(r *http.Request, key string) (*time.Time, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, ErrBadRequest
	}
	return &parsed, nil
}
