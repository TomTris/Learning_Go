package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCreateIncidentRequest_Validate(t *testing.T) {
	valid := func() CreateIncidentRequest {
		return CreateIncidentRequest{
			Title: "outage", Service: "api", Severity: "SEV1",
		}
	}

	t.Run("valid request", func(t *testing.T) {
		r := valid()
		if err := r.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("empty title", func(t *testing.T) {
		r := valid()
		r.Title = ""
		if !errors.Is(r.Validate(), ErrNoTitle) {
			t.Error("expected ErrNoTitle")
		}
	})

	t.Run("whitespace title", func(t *testing.T) {
		r := valid()
		r.Title = "   "
		if !errors.Is(r.Validate(), ErrNoTitle) {
			t.Error("expected ErrNoTitle")
		}
	})

	t.Run("empty service", func(t *testing.T) {
		r := valid()
		r.Service = ""
		if !errors.Is(r.Validate(), ErrNoService) {
			t.Error("expected ErrNoService")
		}
	})

	t.Run("invalid severity", func(t *testing.T) {
		r := valid()
		r.Severity = "SEV4"
		if !errors.Is(r.Validate(), ErrInvalidSeverity) {
			t.Error("expected ErrInvalidSeverity")
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		r := CreateIncidentRequest{
			Title: "  outage  ", Service: "  api  ", Severity: "  SEV1  ",
		}
		r.Validate()
		if r.Title != "outage" {
			t.Errorf("Title not trimmed: %q", r.Title)
		}
		if r.Service != "api" {
			t.Errorf("Service not trimmed: %q", r.Service)
		}
	})
}

func TestTimelineEntry_Validate(t *testing.T) {
	valid := func() TimelineEntry {
		return TimelineEntry{Author: "anh", Type: OBSERVATION, Text: "cpu high"}
	}

	t.Run("valid entry", func(t *testing.T) {
		e := valid()
		if err := e.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	// t.Run("empty author", func(t *testing.T) {
	// 	e := valid()
	// 	e.Author = "   "
	// 	if !errors.Is(e.Validate(), ErrNoAuthor) {
	// 		t.Error("expected ErrNoAuthor")
	// 	}
	// })

	t.Run("invalid type", func(t *testing.T) {
		e := valid()
		e.Type = "active"
		if !errors.Is(e.Validate(), ErrBadEntryType) {
			t.Error("expected ErrBadEntryType")
		}

		e = valid()
		e.Type = "invalid"
		if !errors.Is(e.Validate(), ErrBadEntryType) {
			t.Error("expected ErrBadEntryType")
		}
	})

	t.Run("empty text", func(t *testing.T) {
		e := valid()
		e.Text = ""
		if !errors.Is(e.Validate(), ErrNoText) {
			t.Error("expected ErrNoText")
		}
	})

	t.Run("all valid entry types", func(t *testing.T) {
		for _, typ := range []string{OBSERVATION, ACTION, DISCOVERY, OPEN_QUESTION, STATE_CHANGE} {
			e := valid()
			e.Type = typ
			if err := e.Validate(); err != nil {
				t.Errorf("type %s should be valid, got %v", typ, err)
			}
		}
	})
}

func TestIncidentFilter_Validate(t *testing.T) {
	t.Run("empty filter valid", func(t *testing.T) {
		f := IncidentFilter{}
		if err := f.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("valid status", func(t *testing.T) {
		for _, s := range []string{TRIGGERED, ACKNOWLEDGED, INVESTIGATING, MITIGATED, RESOLVED, "active"} {
			f := IncidentFilter{Status: s}
			if err := f.Validate(); err != nil {
				t.Errorf("status %s should be valid, got %v", s, err)
			}
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		f := IncidentFilter{Status: "abc"}
		if !errors.Is(f.Validate(), ErrBadIncidentStatus) {
			t.Error("expected ErrBadIncidentStatus")
		}
	})

	t.Run("service passes through", func(t *testing.T) {
		f := IncidentFilter{Service: "api"}
		if err := f.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("empty service not fail", func(t *testing.T) {
		f := IncidentFilter{Service: "  "}
		if err := f.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})
}

func TestIncidentUpdate_Validate(t *testing.T) {
	t.Run("valid status", func(t *testing.T) {
		u := IncidentUpdate{Status: new(RESOLVED)}
		if err := u.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		u := IncidentUpdate{Status: new("active")}
		if !errors.Is(u.Validate(), ErrBadIncidentStatus) {
			t.Error("expected ErrBadIncidentStatus")
		}
	})

	t.Run("valid severity", func(t *testing.T) {
		u := IncidentUpdate{Severity: new(SEV2)}
		if err := u.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("invalid severity", func(t *testing.T) {
		u := IncidentUpdate{Severity: new("SEV9")}
		if !errors.Is(u.Validate(), ErrInvalidSeverity) {
			t.Error("expected ErrInvalidSeverity")
		}
	})

	t.Run("valid on_call", func(t *testing.T) {
		u := IncidentUpdate{OnCall: new("bernd")}
		if err := u.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err.Error())
		}
	})

	t.Run("empty on_call", func(t *testing.T) {
		u := IncidentUpdate{OnCall: new("")}
		if !errors.Is(u.Validate(), ErrOnCall) {
			t.Error("expected ErrOnCall")
		}
	})

	t.Run("all fields nil", func(t *testing.T) {
		u := IncidentUpdate{}
		if !errors.Is(u.Validate(), ErrBadRequest) {
			t.Error("expected ErrBadRequest")
		}
	})

	t.Run("set many fields", func(t *testing.T) {
		u := IncidentUpdate{
			Status:   new(RESOLVED),
			Severity: new(SEV2),
			OnCall:   new("bernd1111"),
		}
		u.Validate()
		if *u.Status != RESOLVED {
			t.Errorf("expected RESOLVED, got %v", u.Status)
		}
		if *u.Severity != SEV2 {
			t.Errorf("expected SEV, got %v", u.Severity)
		}
		if *u.OnCall != "bernd1111" {
			t.Errorf("expected bernd1111, get %v", u.OnCall)
		}
	})
}

func TestBuildHandoffBrief(t *testing.T) {
	now := time.Now()

	t.Run("counts actions and open questions", func(t *testing.T) {
		inc := Incident{
			Severity:  SEV1,
			Status:    TRIGGERED,
			Service:   "api",
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: ACTION, Text: "restarted"},
				{Author: "anh", Type: OPEN_QUESTION, Text: "why?"},
				{Author: "anh", Type: OBSERVATION, Text: "cpu high"},
				{Author: "anh", Type: ACTION, Text: "scaled up"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TakenActions != 2 {
			t.Errorf("TakenActions expected 2, got %d", brief.TakenActions)
		}
		if brief.OpenQuestion != 1 {
			t.Errorf("OpenQuestion expected 1, got %d", brief.OpenQuestion)
		}
		if brief.TotalEntry != 4 {
			t.Errorf("TotalEntry expected 4, got %d", brief.TotalEntry)
		}
	})

	t.Run("handoff count tracks author changes", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: OBSERVATION, Text: "a"},
				{Author: "anh", Type: OBSERVATION, Text: "b"},
				{Author: "bernd", Type: OBSERVATION, Text: "c"},
				{Author: "anh", Type: OBSERVATION, Text: "d"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.HandoffCount != 2 {
			t.Errorf("HandoffCount expected 2, got %d", brief.HandoffCount)
		}
	})

	t.Run("single author zero handoffs", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: OBSERVATION, Text: "a"},
				{Author: "anh", Type: OBSERVATION, Text: "b"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.HandoffCount != 0 {
			t.Errorf("HandoffCount expected 0, got %d", brief.HandoffCount)
		}
	})

	t.Run("empty entries", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TotalEntry != 0 {
			t.Errorf("TotalEntry expected 0, got %d", brief.TotalEntry)
		}
		if brief.HandoffCount != 0 {
			t.Errorf("HandoffCount expected 0, got %d", brief.HandoffCount)
		}
	})

	t.Run("maps incident fields correctly", func(t *testing.T) {
		inc := Incident{
			Severity:  SEV2,
			Status:    INVESTIGATING,
			Service:   "payments",
			CreatedAt: now.Add(-30 * time.Minute),
			Entries:   []TimelineEntry{},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.Severity != SEV2 {
			t.Errorf("Severity expected %s, got %s", SEV2, brief.Severity)
		}
		if brief.Status != INVESTIGATING {
			t.Errorf("Status expected %s, got %s", INVESTIGATING, brief.Status)
		}
		if brief.Service != "payments" {
			t.Errorf("Service expected payments, got %s", brief.Service)
		}
		if brief.ElapsedMinute < 29 || brief.ElapsedMinute > 31 {
			t.Errorf("ElapsedMinute expected ~30, got %d", brief.ElapsedMinute)
		}
		if !brief.CreatedAt.Equal(inc.CreatedAt) {
			t.Errorf("CreatedAt mismatch")
		}
	})

	t.Run("nil flagStore skips detailed brief", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
		if brief.OpenQuestionList != nil {
			t.Error("expected nil OpenQuestionList")
		}
	})

	t.Run("detailed brief when flag enabled", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  true,
			Rollout:  100,
			Variants: []string{"detailed"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: ACTION, Text: "restarted"},
				{Author: "anh", Type: OPEN_QUESTION, Text: "why?"},
			},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList == nil {
			t.Fatal("expected non-nil TakenActionsList")
		}
		if len(*brief.TakenActionsList) != 1 {
			t.Errorf("expected 1 action, got %d", len(*brief.TakenActionsList))
		}
		if brief.OpenQuestionList == nil {
			t.Fatal("expected non-nil OpenQuestionList")
		}
		if len(*brief.OpenQuestionList) != 1 {
			t.Errorf("expected 1 question, got %d", len(*brief.OpenQuestionList))
		}
	})

	t.Run("no detailed brief when flag disabled", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  false,
			Rollout:  100,
			Variants: []string{"detailed"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})

	t.Run("no detailed brief when flag not found", func(t *testing.T) {
		fs := CreateFlagStore()
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})

	t.Run("no detailed brief when variant is not detailed", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  true,
			Rollout:  100,
			Variants: []string{"control"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})
}

func TestMarshalNewEntryEvent(t *testing.T) {
	timelineEntry := TimelineEntry{
		ID:        "TLE-1",
		CreatedAt: time.Now(),
		Author:    "anh",
		Type:      OBSERVATION,
		Text:      "test entry",
	}
	rawMsg := marshalNewEntryEvent("INC-test1", timelineEntry)
	var event map[string]any
	json.Unmarshal(rawMsg, &event)

	if event["type"] != "new_entry" {
		t.Fatalf("type expected %v, get %v", OBSERVATION, event["type"])
	}
	if event["incident_id"] != "INC-test1" {
		t.Fatalf("incident_id %v, get %v", "INC-test1", event["incident_id"])
	}
	e := event["entry"].(map[string](any))
	if e["author"] != "anh" {
		t.Fatalf("author expected %v, get %v", "anh", e["author"])
	}
	if e["type"] != OBSERVATION {
		t.Fatalf("type expected %v, get %v", OBSERVATION, e["type"])
	}
}

func TestMarshalIncidentUpdateEvent(t *testing.T) {
	inc := Incident{
		ID:        "INC-test1",
		Title:     "test title",
		Service:   "test service",
		Severity:  "SEV1",
		Status:    TRIGGERED,
		OpenedBy:  "anh",
		OnCall:    "tom",
		CreatedAt: time.Now().Add(-15 * time.Minute),
		UpdatedAt: time.Now(),
		Entries:   []TimelineEntry{},
	}

	rawMsg := marshalIncidentUpdateEvent(inc)
	var event map[string]any
	json.Unmarshal(rawMsg, &event)
	if event["type"] != "incident_updated" {
		t.Fatalf("type expected %v, get %v", "incident_updated", event["type"])
	}
	if event["type"] != "incident_updated" {
		t.Fatalf("type expected %v, get %v", "incident_updated", event["type"])
	}
	e := event["incident"].(map[string]any)
	if e["id"] != "INC-test1" {
		t.Fatalf("id expected %v, get %v", "INC-test1", e["id"])
	}
	if e["service"] != "test service" {
		t.Fatalf("id expected %v, get %v", "test service", e["service"])
	}
}

func TestGetIncidentOK(t *testing.T) {
	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	validIncRequest := validCreateIncidentRequest()
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", validIncRequest)

	handler := IncidentHandler{IncidentStore: MemoryIncidentStore}
	req := httptest.NewRequest("GET", "/incident/INC-1", nil)
	req.SetPathValue("id", "INC-1")

	res, err := handler.GetIncident(req)

	if err != nil {
		t.Fatalf("expected no error, get error %v", err.Error())
	}
	if res.Status != http.StatusOK {
		t.Fatalf("expected status %v, get %v", http.StatusOK, res.Status)
	}
	inc := res.Body.(Incident)
	if inc.ID != "INC-1" {
		t.Fatalf("expected id %v, get %v", "INC-1", inc.ID)
	}
	if inc.Title != validIncRequest.Title {
		t.Fatalf("expected Title %v, get %v", validIncRequest.Title, inc.Title)
	}
	if inc.Service != validIncRequest.Service {
		t.Fatalf("expected Service %v, get %v", validIncRequest.Service, inc.Service)
	}
	if inc.Severity != validIncRequest.Severity {
		t.Fatalf("expected Severity %v, get %v", validIncRequest.Severity, inc.Severity)
	}
	if inc.OnCall != "" {
		t.Fatalf("expected OnCall %v, get %v", "", inc.OnCall)
	}
}

func TestGetIncident404(t *testing.T) {
	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	handler := IncidentHandler{IncidentStore: MemoryIncidentStore}
	req := httptest.NewRequest("GET", "/incident/INC-1", nil)
	req.SetPathValue("id", "INC-1")

	_, appErr := handler.GetIncident(req)

	if appErr == nil {
		t.Fatal("expect error")
	}
	if appErr.Status != 404 {
		t.Fatalf("expected code 404, get %v", appErr.Status)
	}
}

func TestCreateIncident(t *testing.T) {
	incCreateRequest := validCreateIncidentRequest()
	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", incCreateRequest)
	NewInMemoryOnCallStore, _ := NewInMemoryOnCallStore()
	onCallHandler := &OnCallHandler{Store: NewInMemoryOnCallStore}
	handler := IncidentHandler{
		IncidentStore: MemoryIncidentStore,
		CurrentOnCall: onCallHandler.Store,
	}
	bodyRaw, _ := json.Marshal(incCreateRequest)

	req := httptest.NewRequest("POST", "/incident", bytes.NewReader(bodyRaw))
	ctx := context.WithValue(req.Context(), userContextKey, UserContext{
		ID:       "Usr-1",
		Username: "username_admin",
		Role:     "admin",
	})
	appRes, err := handler.CreateIncident(req.WithContext(ctx))

	if err != nil {
		t.Fatalf("expected no error, get %v", err.Error())
	}
	if appRes.Status != http.StatusCreated {
		t.Fatalf("status code expected %v, get %v", http.StatusCreated, appRes.Status)
	}
	// Evaluate
	response := appRes.Body.(Incident)
	if response.ID != "INC-2" {
		t.Fatalf("status code expected %v, get %v", "INC-2", response.ID)
	}
	if response.Title != incCreateRequest.Title {
		t.Fatalf("title expected %v, got %v", incCreateRequest.Title, response.Title)
	}
	if response.Severity != incCreateRequest.Severity {
		t.Fatalf("Severity expected %v, got %v", incCreateRequest.Severity, response.Severity)
	}
	if response.Service != incCreateRequest.Service {
		t.Fatalf("Service expected %v, got %v", incCreateRequest.Service, response.Service)
	}
	if response.OpenedBy != "username_admin" {
		t.Fatalf("OpenedBy expected %v, got %v", "username_admin", response.OpenedBy)
	}
}

func TestListIncident(t *testing.T) {

	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	incCreateRequest := validCreateIncidentRequest()
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", incCreateRequest)
	incCreateRequest.Title = "123"
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", incCreateRequest)
	incCreateRequest.Service = "no_services"
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", incCreateRequest)
	incCreateRequest.Severity = "SEV3"
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", incCreateRequest)

	handler := IncidentHandler{IncidentStore: MemoryIncidentStore}

	t.Run("listAll", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents", nil)
		appRes, err := handler.ListIncidents(req)
		if err != nil {
			t.Fatalf("expected no error, get %v", err.Error())
		}
		if appRes.Status != http.StatusOK {
			t.Fatalf("status code expected %v, get %v", http.StatusOK, appRes.Status)
		}

		// Evaluate
		response := appRes.Body.([]Incident)
		if len(response) != 4 {
			t.Fatalf("len expect %v, get %v", 4, len(response))
		}
	})
	t.Run("listByStatus", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents?service=no_services", nil)
		appRes, err := handler.ListIncidents(req)
		if err != nil {
			t.Fatalf("expected no error, get %v", err.Error())
		}
		if appRes.Status != http.StatusOK {
			t.Fatalf("status code expected %v, get %v", http.StatusOK, appRes.Status)
		}

		// Evaluate
		response := appRes.Body.([]Incident)
		if len(response) != 2 {
			t.Fatalf("len expect %v, get %v", 2, len(response))
		}
	})
}

// init a server with an incident available
func newTestServer(t *testing.T) (*httptest.Server, string, string) {
	t.Helper()

	promRegistry := prometheus.NewRegistry()
	httpMetrics := NewHttpMetrics(promRegistry)
	metricRegistry := NewMetricRegistry(promRegistry)
	incidentStoreMetric := NewIncidentStoreMetric(promRegistry)

	registry := NewRegistry(metricRegistry)
	go registry.run()
	t.Cleanup(func() { close(registry.done) })

	NewInMemoryOnCallStore, _ := NewInMemoryOnCallStore()
	onCallHandler := &OnCallHandler{Store: NewInMemoryOnCallStore}
	flagHandler := FlagHandler{store: CreateFlagStore()}
	memStore, _ := NewMemoryIncidentStore()
	instrumentedIncidentStore := InstrumentedIncidentStore{
		inner:   memStore,
		metrics: incidentStoreMetric,
	}
	incHandler := IncidentHandler{
		IncidentStore: &instrumentedIncidentStore,
		Registry:      registry,
		FlagEvaluator: &flagHandler.store,
		CurrentOnCall: onCallHandler.Store,
	}

	pwd1, err1 := HashPassword("anh123")
	pwd2, err2 := HashPassword("bernd123")
	pwd3, err3 := HashPassword("admin123")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("HashPassword has problem")
	}

	var seedUsers = []User{
		{ID: "u1", Username: "anh", Password: pwd1, Role: "engineer"},
		{ID: "u2", Username: "bernd", Password: pwd2, Role: "engineer"},
		{ID: "u3", Username: "admin", Password: pwd3, Role: "admin"},
	}
	userStore := NewMemoryUserStoreWithSeed(seedUsers)
	jwt_secret := "testing-JWT-secret"
	authHandler := NewAuthHandler(userStore, []byte(jwt_secret), time.Duration(15))
	ttl := time.Duration(15 * time.Minute)
	engineerTokenSigned, _ := IssueToken(seedUsers[0], []byte(jwt_secret), ttl, time.Now())
	adminTokenSigned, _ := IssueToken(seedUsers[2], []byte(jwt_secret), ttl, time.Now())

	memStore.CreateIncident(context.Background(), "", "", validCreateIncidentRequest())

	router := getRouter(&incHandler, &flagHandler, authHandler, onCallHandler, nil, promRegistry, httpMetrics)
	return httptest.NewServer(router), engineerTokenSigned, adminTokenSigned
}

func TestAddEntry(t *testing.T) {
	srv, engineerTokenSigned, adminTokenSigned := newTestServer(t)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/incidents/INC-1/ws"
	jar, _ := cookiejar.New(nil)
	srvURL, _ := url.Parse(srv.URL)
	jar.SetCookies(srvURL, []*http.Cookie{{Name: "access_token", Value: engineerTokenSigned}})
	dialer := websocket.Dialer{Jar: jar}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("expected no error, get error %v", err.Error())
	}
	defer conn.Close()

	entry := TimelineEntry{
		Type: OBSERVATION,
		Text: "looking into A",
	}
	bodyRaw, _ := json.Marshal(entry)

	// HTTP Respsone
	t.Run("Test normal request with admin role", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/incidents/INC-1/entries", bytes.NewReader(bodyRaw))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "access_token", Value: adminTokenSigned})
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("status code expected %v, got %v", http.StatusCreated, resp.StatusCode)
		}

		var resEntry1 TimelineEntry
		json.NewDecoder(resp.Body).Decode(&resEntry1)

		if resEntry1.ID != "TLE-1" {
			t.Fatalf("entry ID expected %v, got %v", "TLE-1", resEntry1.ID)
		}
		if resEntry1.Type != OBSERVATION {
			t.Fatalf("entry type expected %v, got %v", OBSERVATION, resEntry1.Type)
		}

		// Websocket Response
		_, msgRaw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Expected no error, get %v", err)
		}
		var wsMsg map[string]any
		json.Unmarshal(msgRaw, &wsMsg)
		if wsMsg["type"] != "new_entry" {
			t.Fatalf("expected type %v, get %v", "new_entry", wsMsg["type"])
		}
		if wsMsg["incident_id"] != "INC-1" {
			t.Fatalf("expected Incident ID %v, get %v", "INC-1", wsMsg["incident_id"])
		}

		entryRaw, err := json.Marshal(wsMsg["entry"])
		if err != nil {
			t.Fatalf("Expected no error, got %v", err.Error())
		}

		var respEntry2 map[string]string
		json.Unmarshal(entryRaw, &respEntry2)

		if respEntry2["author"] != "admin" {
			t.Fatalf("author expected %v, get %v", "admin", respEntry2["author"])
		}
		if respEntry2["type"] != OBSERVATION {
			t.Fatalf("type expected %v, get %v", OBSERVATION, respEntry2["type"])
		}
	})

	t.Run("Test normal request with engineer role", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/incidents/INC-1/entries", bytes.NewReader(bodyRaw))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "access_token", Value: engineerTokenSigned})
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("status code expected %v, got %v", http.StatusCreated, resp.StatusCode)
		}
	})

}

func TestUpdateIncident(t *testing.T) {
	srv, engineerTokenSigned, adminTokenSigned := newTestServer(t)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/incidents/INC-1/ws"
	jar, _ := cookiejar.New(nil)
	srvURL, _ := url.Parse(srv.URL)
	jar.SetCookies(srvURL, []*http.Cookie{{Name: "access_token", Value: engineerTokenSigned}})
	dialer := websocket.Dialer{Jar: jar}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("expected no error, get error %v", err.Error())
	}
	defer conn.Close()

	incidentUpdate := IncidentUpdate{
		Status:   new(RESOLVED),
		Severity: new("SEV2"),
	}
	bodyRaw, _ := json.Marshal(incidentUpdate)

	// HTTP Respsone
	t.Run("test with engineer role", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/api/incidents/INC-1", bytes.NewReader(bodyRaw))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "access_token", Value: engineerTokenSigned})
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status code expected %v, got %v", http.StatusForbidden, resp.StatusCode)
		}
	})
	t.Run("test with admin role", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/api/incidents/INC-1", bytes.NewReader(bodyRaw))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "access_token", Value: adminTokenSigned})
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("status code expected %v, got %v", http.StatusNoContent, resp.StatusCode)
		}

		// Websocket Response
		_, msgRaw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Expected no error, get %v", err)
		}
		var wsMsg map[string]any
		json.Unmarshal(msgRaw, &wsMsg)
		if wsMsg["type"] != "incident_updated" {
			t.Fatalf("expected type %v, get %v", "incident_updated", wsMsg["type"])
		}

		incidentRaw, err := json.Marshal(wsMsg["incident"])
		if err != nil {
			t.Fatalf("Expected no error, got %v", err.Error())
		}

		var resIncident map[string]string
		json.Unmarshal(incidentRaw, &resIncident)

		if resIncident["id"] != "INC-1" {
			t.Fatalf("expected Incident ID %v, get %v", "INC-1", wsMsg["incident_id"])
		}
		if resIncident["status"] != RESOLVED {
			t.Fatalf("author expected %v, get %v", RESOLVED, resIncident["status"])
		}
		if resIncident["severity"] != "SEV2" {
			t.Fatalf("type expected %v, get %v", "SEV2", resIncident["severity"])
		}
	})
}

func TestGetHandoffBriefAdmin(t *testing.T) {
	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	validIncRequest := validCreateIncidentRequest()
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", validIncRequest)

	handler := IncidentHandler{IncidentStore: MemoryIncidentStore}
	req := httptest.NewRequest("GET", "/incidents/INC-1/handoff?user_id=tom", nil)
	req.SetPathValue("id", "INC-1")
	ctx := context.WithValue(req.Context(), userContextKey, UserContext{
		ID:       "Usr-1",
		Username: "username_admin",
		Role:     "admin",
	})

	appRes, appErr := handler.GetHandoffBrief(req.WithContext(ctx))
	if appErr != nil {
		t.Fatalf("expected no error, get %v", appErr.Error())
	}
	if appRes.Status != http.StatusOK {
		t.Fatalf("expected status %v, get %v", http.StatusOK, appRes.Status)
	}

	bodyRaw, err := json.Marshal(appRes.Body)
	if err != nil {
		t.Fatalf("expected not nil, get %v", err.Error())
	}
	var body HandoffBrief
	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		t.Fatalf("expected not nil, get %v", err.Error())
	}
}

func TestGetHandoffBriefEngineer(t *testing.T) {
	MemoryIncidentStore, _ := NewMemoryIncidentStore()
	validIncRequest := validCreateIncidentRequest()
	MemoryIncidentStore.CreateIncident(context.Background(), "", "", validIncRequest)

	handler := IncidentHandler{IncidentStore: MemoryIncidentStore}
	req := httptest.NewRequest("GET", "/incidents/INC-1/handoff", nil)
	req.SetPathValue("id", "INC-1")
	ctx := context.WithValue(req.Context(), userContextKey, UserContext{
		ID:       "Usr-1",
		Username: "username_engineer",
		Role:     "engineer",
	})

	appRes, appErr := handler.GetHandoffBrief(req.WithContext(ctx))
	if appErr != nil {
		t.Fatalf("expected no error, get %v", appErr.Error())
	}
	if appRes.Status != http.StatusOK {
		t.Fatalf("expected status %v, get %v", http.StatusOK, appRes.Status)
	}

	bodyRaw, err := json.Marshal(appRes.Body)
	if err != nil {
		t.Fatalf("expected not nil, get %v", err.Error())
	}
	var body HandoffBrief
	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		t.Fatalf("expected not nil, get %v", err.Error())
	}
}
