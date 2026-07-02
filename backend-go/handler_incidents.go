package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type FlagEvaluator interface {
	Evaluate(flagName string, userID string) (*FlagEvaluateAnswer, error)
}

type IncidentHandler struct {
	IncidentStore IncidentStore
	Registry      Registry
	FlagEvaluator FlagEvaluator
	CurrentOnCall CurrentOnCall
}

type Incident struct {
	ID        string          `json:"id" bson:"_id,omitempty"`
	Title     string          `json:"title" bson:"title"`
	Service   string          `json:"service" bson:"service"`
	Severity  string          `json:"severity" bson:"severity"` // SEV1, SEV2, SEV3
	Status    string          `json:"status" bson:"status"`     // triggered, acknowledged, investigating, mitigated, resolved
	OpenedBy  string          `json:"opened_by" bson:"opened_by"`
	OnCall    string          `json:"on_call" bson:"on_call"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
	Entries   []TimelineEntry `json:"entries" bson:"entries"`
	Version   int             `json:"version" bson:"version"`
}

type CreateIncidentRequest struct {
	Title    string `json:"title" bson:"title"`
	Service  string `json:"service" bson:"service"`
	Severity string `json:"severity" bson:"severity"` // SEV1, SEV2, SEV3
}

func (c *CreateIncidentRequest) Validate() error {
	c.Title = strings.TrimSpace(c.Title)
	if c.Title == "" {
		return ErrNoTitle
	}

	c.Service = strings.TrimSpace(c.Service)
	if c.Service == "" {
		return ErrNoService
	}

	c.Severity = strings.TrimSpace(c.Severity)
	if IncidentSeverity[c.Severity] == false {
		return ErrInvalidSeverity
	}
	return nil
}

type TimelineEntry struct {
	ID        string    `json:"id" bson:"id"`
	Author    string    `json:"author" bson:"author"`
	Type      string    `json:"type" bson:"type"` // observation, action, discovery, open_question, state_change
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (c *TimelineEntry) Validate() error {
	c.Type = strings.TrimSpace(c.Type)
	if validEntryTypes[c.Type] == false {
		return ErrBadEntryType
	}
	c.Text = strings.TrimSpace(c.Text)
	if c.Text == "" {
		return ErrNoText
	}
	return nil
}

type IncidentFilter struct {
	Status  string `json:"status" bson:"status"`
	Service string `json:"service" bson:"service"`
}

func (f *IncidentFilter) Validate() error {
	f.Status = strings.TrimSpace(f.Status)
	f.Service = strings.TrimSpace(f.Service)
	if f.Status != "" && !IncidentStatus[f.Status] {
		return ErrBadIncidentStatus
	}
	return nil
}

type IncidentUpdate struct {
	Status   *string `json:"status,omitempty" bson:"status"`
	Severity *string `json:"severity,omitempty" bson:"severity"`
	OnCall   *string `json:"on_call,omitempty" bson:"on_call"`
}

func (f *IncidentUpdate) Validate() error {
	if f.Status == nil && f.Severity == nil && f.OnCall == nil {
		return ErrBadRequest
	}

	if f.Status != nil {
		trimmed := strings.TrimSpace(*f.Status)
		if IncidentStatus[trimmed] == false || trimmed == "active" {
			return ErrBadIncidentStatus
		}
		f.Status = &trimmed
	}
	if f.Severity != nil {
		trimmed := strings.TrimSpace(*f.Severity)
		if IncidentSeverity[trimmed] == false {
			return ErrInvalidSeverity
		}
		f.Severity = &trimmed
	}
	if f.OnCall != nil {
		trimmed := strings.TrimSpace(*f.OnCall)
		if trimmed == "" {
			return ErrOnCall
		}
		f.OnCall = &trimmed
	}
	return nil
}

type HandoffBrief struct {
	Severity         string           `json:"severity"`
	Status           string           `json:"status"`
	Service          string           `json:"service"`
	TotalEntry       int              `json:"total_entry"`
	ElapsedMinute    int              `json:"elapsed_minute"`
	TakenActions     int              `json:"taken_actions"`
	OpenQuestion     int              `json:"open_question"`
	HandoffCount     int              `json:"handoff_count"`
	TakenActionsList *[]TimelineEntry `json:"taken_actions_list,omitempty"`
	OpenQuestionList *[]TimelineEntry `json:"open_question_list,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}

func buildHandoffBrief(inc Incident, flagEvaluator FlagEvaluator, userID string) HandoffBrief {
	actionsList := []TimelineEntry{}
	openQuestionsList := []TimelineEntry{}
	author := ""
	handoffCount := 0
	for _, entry := range inc.Entries {
		if author != entry.Author {
			author = entry.Author
			handoffCount++
		}
		switch entry.Type {
		case ACTION:
			actionsList = append(actionsList, entry)
		case OPEN_QUESTION:
			openQuestionsList = append(openQuestionsList, entry)
		}
	}

	if handoffCount != 0 {
		handoffCount--
	}

	brief := HandoffBrief{
		Severity:      inc.Severity,
		Status:        inc.Status,
		Service:       inc.Service,
		ElapsedMinute: int(time.Since(inc.CreatedAt).Minutes()),
		TotalEntry:    len(inc.Entries),
		TakenActions:  len(actionsList),
		OpenQuestion:  len(openQuestionsList),
		HandoffCount:  handoffCount,
		CreatedAt:     inc.CreatedAt,
	}

	if flagEvaluator != nil {
		flagAnswer, err := flagEvaluator.Evaluate("detailed_handoff_brief", userID)
		if err == nil && flagAnswer.InRollout == true && *flagAnswer.Variant == "detailed" {
			brief.OpenQuestionList = &openQuestionsList
			brief.TakenActionsList = &actionsList
		}
	}
	return brief
}

func marshalNewEntryEvent(incidentID string, entry TimelineEntry) json.RawMessage {
	event := struct {
		Type       string        `json:"type"`
		IncidentID string        `json:"incident_id"`
		Entry      TimelineEntry `json:"entry"`
	}{
		Type:       "new_entry",
		IncidentID: incidentID,
		Entry:      entry,
	}
	data, _ := json.Marshal(event)
	return data
}

func marshalIncidentUpdateEvent(incAfter Incident) json.RawMessage {
	event := struct {
		Type     string   `json:"type"`
		Incident Incident `json:"incident"`
	}{
		Type:     "incident_updated",
		Incident: incAfter,
	}

	data, _ := json.Marshal(event)
	return data
}

func (h *IncidentHandler) CreateIncident(r *http.Request) (*AppResponse, *AppError) {
	req := CreateIncidentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}

	if err := req.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	onCall, err := h.CurrentOnCall.CurrentOnCall(r.Context(), req.Service)
	if err != nil {
		if errors.Is(err, OnCallShiftEntryNotFound) {
			onCall = ""
		} else {
			return nil, InternalServerError(err)
		}
	}

	user := r.Context().Value(userContextKey).(UserContext)
	createdIncident, err := h.IncidentStore.CreateIncident(r.Context(), user.Username, onCall, req)
	if err != nil {
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusCreated, createdIncident), nil
}

func (h *IncidentHandler) GetIncident(r *http.Request) (*AppResponse, *AppError) {
	incidentID := r.PathValue("id")
	inc, err := h.IncidentStore.GetIncident(r.Context(), incidentID)

	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	return newAppResponse(http.StatusOK, inc), nil
}

// TODO Handle error cases properly
func (h *IncidentHandler) AddEntry(r *http.Request) (*AppResponse, *AppError) {
	user := r.Context().Value(userContextKey).(UserContext)
	inc, err := h.IncidentStore.GetIncident(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}
	if err := AuthorizeIncidentAction(user, inc, ActionAddEntry); err != nil {
		return nil, Forbidden(err)
	}

	timelineEntry := TimelineEntry{}
	if err := json.NewDecoder(r.Body).Decode(&timelineEntry); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}
	if err := timelineEntry.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	timelineEntry.Author = user.Username
	newEntry, err := h.IncidentStore.AddEntry(r.Context(), inc.ID, inc.Version, timelineEntry)

	if err != nil {
		if errors.Is(err, ErrIncidentVersionConflict) || errors.Is(err, ErrIncidentResolved) {
			return nil, Conflict(err)
		}
		return nil, InternalServerError(err)
	}

	h.Registry.broadcast <- BroadcastMessage{
		incidentID: inc.ID,
		msg:        marshalNewEntryEvent(inc.ID, newEntry),
	}

	return newAppResponse(http.StatusCreated, newEntry), nil
}

func (h *IncidentHandler) ListIncidents(r *http.Request) (*AppResponse, *AppError) {
	incidentFilter := IncidentFilter{
		Status:  r.URL.Query().Get("status"),
		Service: r.URL.Query().Get("service"),
	}

	if err := incidentFilter.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	filteredIncidents, err := h.IncidentStore.ListIncidents(r.Context(), incidentFilter)
	if err != nil {
		return nil, InternalServerError(err)
	}
	return newAppResponse(http.StatusOK, filteredIncidents), nil
}

func (h *IncidentHandler) UpdateIncident(r *http.Request) (*AppResponse, *AppError) {
	user := r.Context().Value(userContextKey).(UserContext)
	inc, err := h.IncidentStore.GetIncident(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}
	if err := AuthorizeIncidentAction(user, inc, ActionUpdateIncident); err != nil {
		return nil, Forbidden(err)
	}

	incidentUpdate := IncidentUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&incidentUpdate); err != nil {
		return nil, BadRequest(MalformedRequestBody)
	}
	if err := incidentUpdate.Validate(); err != nil {
		return nil, BadRequest(err)
	}

	incAfter, err := h.IncidentStore.UpdateIncident(r.Context(), inc.ID, inc.Version, incidentUpdate)
	if err != nil {
		if errors.Is(err, ErrIncidentVersionConflict) {
			return nil, Conflict(err)
		}
		return nil, InternalServerError(err)
	}

	h.Registry.broadcast <- BroadcastMessage{
		msg:        marshalIncidentUpdateEvent(incAfter),
		incidentID: inc.ID,
	}
	return newAppResponse(http.StatusNoContent, nil), nil
}

func (h *IncidentHandler) GetHandoffBrief(r *http.Request) (*AppResponse, *AppError) {
	incidentID := r.PathValue("id")
	inc, err := h.IncidentStore.GetIncident(r.Context(), incidentID)

	if err != nil {
		if errors.Is(err, ErrIncidentNotFound) {
			return nil, NotFound(err)
		}
		return nil, InternalServerError(err)
	}

	user := r.Context().Value(userContextKey).(UserContext)
	var body HandoffBrief
	if user.Role == "admin" {
		userID := r.URL.Query().Get("user_id")
		body = buildHandoffBrief(inc, h.FlagEvaluator, userID)
	} else {
		body = buildHandoffBrief(inc, h.FlagEvaluator, user.ID)
	}
	return newAppResponse(http.StatusOK, body), nil
}

func (h *IncidentHandler) HandleIncidentWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		RequestID := getRequestID(r)
		writeError(w, http.StatusInternalServerError, ErrorMessageJSON{
			ErrorCode: INTERNAL_SERVER_ERROR,
			Message:   err.Error(),
			RequestID: RequestID,
		})
		return
	}

	incidentID := r.PathValue("id")
	client := newClient(incidentID, conn)
	client.joinRegistry(&h.Registry)

	go client.writePump()
	go client.readPump()
}
