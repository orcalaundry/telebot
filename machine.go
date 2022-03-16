package main

import (
	"fmt"
	"strings"
	"time"
)

// machine represents either a washer or dryer.
type machine struct {
	Name          string        `json:"name"`
	Floor         int           `json:"floor"`
	Pos           int           `json:"pos"`
	Type          machineType   `json:"type"`
	Status        machineStatus `json:"status"`
	LastStartedAt time.Time     `json:"last_started_at"`
	TimeLeft      time.Duration `json:"time_left"`
	Duration      int           `json:"duration"`
}

type machineType string

const (
	washerMachine  machineType = "washer"
	dryerMachine               = "dryer"
	unknownMachine             = "unknown"
)

func (m *machineType) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	switch s {
	case "washer":
		*m = washerMachine
	case "dryer":
		*m = dryerMachine
	default:
		*m = unknownMachine
	}
	return nil
}

// machineStatus serves as an enum type for tracking the machine's state,
// e.g. whether it is idle or in use.
type machineStatus int

const (
	statusError machineStatus = iota
	statusIdle
	statusFinishing
	statusInUse
	statusUnknown
)

// The string representation of a machine's status should complete the sentence
// "This machine is ...".
func (m *machineStatus) String() string {
	switch *m {
	case statusError:
		return "facing an error"
	case statusIdle:
		return "idle"
	case statusFinishing:
		return "almost done"
	case statusInUse:
		return "in use"
	default:
		return "unreachable"
	}
}

func (m *machineStatus) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	switch s {
	case "error":
		*m = statusError
	case "idle":
		*m = statusIdle
	case "finishing":
		*m = statusFinishing
	case "in_use":
		*m = statusInUse
	default:
		*m = statusUnknown
	}
	return nil
}

// machineHumanFriendlyNames maps the position of each machine in the laundry
// room to their conventional names.
var machineHumanFriendlyNames map[int]string = map[int]string{
	0: "Coin washer",
	1: "PayLah washer",
	2: "PayLah dryer",
	3: "Coin dryer",
}

// addName decides on the right resident-friendly name for the machine based
// on the position.
func (m *machine) addName() {
	name, ok := machineHumanFriendlyNames[m.Pos]
	if !ok {
		m.Name = "Unknown machine"
		return
	}
	m.Name = name
}

// computeTimeLeft sets the TimeLeft field on the machine, if the machine is
// in use.
func (m *machine) computeTimeLeft() {
	if m.Status != statusInUse {
		return
	}
	m.TimeLeft = time.Until(m.LastStartedAt.Add(time.Duration(m.Duration) * time.Second))
}

func (m *machine) String() string {
	return fmt.Sprintf("%s is %s \n\t Last started at %s \n\t %s remaining",
		m.Name,
		m.Status.String(),
		m.LastStartedAt.Format(time.Kitchen),
		m.TimeLeft)
}

// machines is a wrapper around a slice of machines.
type machines []*machine

func (ms machines) String() string {
	b := strings.Builder{}
	for _, m := range ms {
		b.WriteString(m.String())
		b.WriteString("\n\n")
	}
	return b.String()
}
