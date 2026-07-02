package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type FlagHandler struct {
	store FlagStore
}

type FeatureFlag struct {
	Name     string   `json:"name"`
	Enabled  bool     `json:"enabled"`
	Rollout  int      `json:"rollout"`  // 0–100, percentage of users who see the feature
	Variants []string `json:"variants"` // e.g., ["control", "variant_a", "variant_b"]
}

func (f *FeatureFlag) Validate() error {
	if strings.TrimSpace(f.Name) == "" {
		return ErrBadRequest
	}
	if f.Rollout < 0 || f.Rollout > 100 {
		return ErrBadRequest
	}
	if len(f.Variants) == 0 {
		return ErrBadRequest
	}
	variants := make(map[string]bool)
	for _, variant := range f.Variants {
		if strings.TrimSpace(variant) == "" || variants[variant] == true {
			return ErrBadRequest
		}
		variants[variant] = true
	}
	return nil
}

type FeatureFlagUpdate struct {
	Name    string `json:"name"`
	Enabled *bool  `json:"enabled"`
	Rollout *int   `json:"rollout"` // 0–100, percentage of users who see the feature
}

func (u *FeatureFlagUpdate) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("Bad Flag Name")
	}
	if u.Enabled == nil && u.Rollout == nil {
		return errors.New("both enabled and rollout are empty")
	}
	if u.Rollout != nil && (*u.Rollout < 0 || *u.Rollout > 100) {
		return errors.New("invalid rollout")
	}
	return nil
}

type FlagEvaluateAnswer struct {
	Name      string  `json:"name"`
	UserID    string  `json:"user_id"`
	Enabled   bool    `json:"enabled"`
	InRollout bool    `json:"in_rollout"`
	Variant   *string `json:"variants"`
}

func (flagHandler *FlagHandler) CreateFlag(r *http.Request) (*AppResponse, *AppError) {
	f := FeatureFlag{}
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}
	if err := f.Validate(); err != nil {
		return nil, BadRequest(err)
	}
	if err := flagHandler.store.Create(f); err != nil {
		if errors.Is(err, ErrFlagAlreadyExist) {
			return nil, Conflict(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusCreated, f), nil
}

func (flagHandler *FlagHandler) UpdateFlag(r *http.Request) (*AppResponse, *AppError) {
	u := FeatureFlagUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}
	u.Name = r.PathValue("name")
	if err := u.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	if err := flagHandler.store.Update(u); err != nil {
		if errors.Is(err, ErrFlagNotfound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusNoContent, nil), nil
}

func (flagHandler *FlagHandler) ListAllFlag(r *http.Request) (*AppResponse, *AppError) {
	allFlags, err := flagHandler.store.AllFlags()
	if err != nil {
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusOK, allFlags), nil
}

func (flagHandler *FlagHandler) Evaluate(r *http.Request) (*AppResponse, *AppError) {
	flagName := r.PathValue("name")
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		return nil, BadRequest(errors.New("empty user_id"))
	}

	flagAnswer, err := flagHandler.store.Evaluate(flagName, userID)
	if err != nil {
		if errors.Is(err, ErrFlagNotfound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusOK, flagAnswer), nil
}
