package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// machine represents either a washer or dryer.
type machine struct {
	Name          string        `json:"name"`
	Floor         int           `json:"floor"`
	Type          machineType   `json:"type"`
	Status        machineStatus `json:"status"`
	LastStartedAt time.Time     `json:"last_started_at"`
	TimeLeft      time.Duration `json:"time_left"`
}

type machineType string

const (
	washer machineType = "washer"
	dryer  machineType = "dryer"
)

// machineStatus serves as an enum type for tracking the machine's state,
// e.g. whether it is idle or in use.
type machineStatus string

const (
	statusError     machineStatus = "error"
	statusIdle      machineStatus = "idle"
	statusFinishing machineStatus = "finishing"
	statusInUse     machineStatus = "in_use"
)

func (m *machine) UnmarshalJSON(data []byte) error {
	// First, we unmarshal into a map. Then we pick out the fields we need.
	var xs map[string]interface{}
	if err := json.Unmarshal(data, &xs); err != nil {
		return err
	}

	// Here, a hot mess of reflection to get all the fields.
	floor_, _ := xs["floor"].(float64)
	floor := int(floor_)
	pos_, _ := xs["pos"].(float64)
	pos := int(pos_)
	typ_, _ := xs["type"].(string)
	typ := machineType(typ_)
	status_, _ := xs["status"].(string)
	status := machineStatus(status_)

	lastStartedAt_, _ := xs["last_started_at"].(string)
	LastStartedAt, err := time.Parse(time.RFC3339, lastStartedAt_)
	if err != nil {
		return err
	}

	// Does this count as magic strings?
	var name string
	switch pos {
	case 0:
		name = "Coin washer"
	case 1:
		name = "PayLah washer"
	case 2:
		name = "PayLah dryer"
	case 3:
		name = "Coin dryer"
	default:
		name = "Unknown machine"
	}

	m.Floor = floor
	m.Type = typ
	m.Status = status
	m.Name = name
	m.LastStartedAt = LastStartedAt

	if status != statusInUse {
		m.TimeLeft = 0
		return nil
	}

	// Perhaps we shouldn't be doing this here. It makes this method impossible
	// to test.
	duration, _ := xs["duration"].(int)
	m.TimeLeft = time.Until(LastStartedAt.Add(time.Duration(duration) * time.Second))

	return nil
}

// machines is a wrapper around a slice of machines.
type machines []*machine

func (ms machines) String() string {
	b := strings.Builder{}
	for _, m := range ms {
		b.WriteString(fmt.Sprintf("%s is %s \n\t last started at %s \n\t %s remaining \n", m.Name, m.Status, m.LastStartedAt.Format(time.Kitchen), m.TimeLeft))
	}
	return b.String()
}
