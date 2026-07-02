package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Config struct {
	Port             string // default "8080",   env: HANDOFF_PORT
	Environment      string // default "development", env: HANDOFF_ENV
	ConnectionString string // default "", env HANDOFF_CONNECT_STRING
	DatabaseName     string
	JWT_SECRET       string
}

func loadConfig() Config {
	if os.Getenv("HANDOFF_ENV") != "production" {
		godotenv.Load("../.env")
	}

	config := Config{
		Port:             envOr("HANDOFF_PORT", "8080"),
		Environment:      envOr("HANDOFF_ENV", "development"),
		ConnectionString: envOr("HANDOFF_CONNECT_STRING", ""),
		DatabaseName:     envOr("HANDOFF_DB", "incident_tracker"),
		JWT_SECRET:       envOr("HANDOFF_JWT_SECRET", ""),
	}
	if len(config.JWT_SECRET) == 0 {
		log.Fatalln("HANDOFF_JWT_SECRET empty")
	}
	return config
}

func envOr(envKey string, defaultValue string) string {
	envValue := os.Getenv(envKey)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	timeout = time.Duration(5 * time.Second)
)

const (
	CollectionUsers        = "users"
	CollectionIncidents    = "incidents"
	CollectionOnCallShifts = "on_call_shifts"

	CollectionCounters              = "counters"
	CollectionCountersUser          = "counter_user"
	CollectionCountersTimelineEntry = "counter_entry_timeline"
	CollectionCountersIncident      = "counter_incident"
	CollectionCountersOnCallShift   = "counter_on_cal_shift"

	requestIDKey   = "request_id"
	userContextKey = "user"

	incidentIDPrefix         = "INC-"
	TimelineEntryIDPrefix    = "TLE-"
	OnCallShiftEntryIDPrefix = "ONc-"
	UserIDPrefix             = "Usr-"

	ARealBrcyptHasedPassword = "$2a$12$gpwXpFXAYxo8U3RoKwYQR.CRfVf6zJzH0TcbR.L0rSrRLD9S00d/e"
)

// Incident Severity
const (
	SEV1 = "SEV1"
	SEV2 = "SEV2"
	SEV3 = "SEV3"
)

var IncidentSeverity = map[string]bool{
	SEV1: true,
	SEV2: true,
	SEV3: true,
}

// Incident status
const (
	TRIGGERED     = "triggered"
	ACKNOWLEDGED  = "acknowledged"
	INVESTIGATING = "investigating"
	MITIGATED     = "mitigated"
	RESOLVED      = "resolved"
)

var IncidentStatus = map[string]bool{
	TRIGGERED:     true,
	ACKNOWLEDGED:  true,
	INVESTIGATING: true,
	MITIGATED:     true,
	RESOLVED:      true,
	"active":      true,
}

// Entry type
const (
	OBSERVATION   = "observation"
	ACTION        = "action"
	DISCOVERY     = "discovery"
	OPEN_QUESTION = "open_question"
	STATE_CHANGE  = "state_change"
)

var validEntryTypes = map[string]bool{
	OBSERVATION:   true,
	ACTION:        true,
	DISCOVERY:     true,
	OPEN_QUESTION: true,
	STATE_CHANGE:  true,
}
