package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_MachineUnmarshalJSON(t *testing.T) {
	tests := []struct {
		json string   // the json string to unmarshal
		want *machine // machine value expected, without the lastStartedAt field
		// lastStartedAtRFC3339 is a formatted string of the machine.LastStartedAt
		// time field, which we will patch into our target machine when the tests
		// run.
		lastStartedAtRFC3339 string
	}{
		{
			json: `{
  "floor": 14,
  "pos": 0,
  "type": "washer",
  "status": "idle",
  "last_started_at": "1970-01-01T00:00:00+00:00",
  "duration": 1800
}`,
			want: &machine{
				Floor:    14,
				Pos:      0,
				Type:     washerMachine,
				Status:   statusIdle,
				Duration: 1800,
			},
			lastStartedAtRFC3339: "1970-01-01T00:00:00+00:00",
		},
		{
			json: `{
  "floor": 14,
  "pos": 2,
  "type": "dryer",
  "status": "in_use",
  "last_started_at": "2022-03-16T05:13:06+00:00",
  "duration": 2400
}`,
			want: &machine{
				Floor:    14,
				Pos:      2,
				Type:     dryerMachine,
				Status:   statusInUse,
				Duration: 2400,
			},
			lastStartedAtRFC3339: "2022-03-16T05:13:06+00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.json, func(t *testing.T) {
			var m machine
			if err := json.Unmarshal([]byte(tt.json), &m); err != nil {
				t.Error(err)
			}
			// Note that tt.want currently contains the zero value for its
			// lastStartedAt field. We need to patch that in.
			lastStartedAt, err := time.Parse(time.RFC3339, tt.lastStartedAtRFC3339)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.LastStartedAt = lastStartedAt
			if !(reflect.DeepEqual(&m, tt.want)) {
				t.Errorf("got %+v, want %+v", &m, tt.want)
			}
		})
	}
}

func Test_MachineComputeTimeLeft(t *testing.T) {
	tests := []struct {
		m    *machine
		want time.Duration
	}{
		{
			m: &machine{
				Status:        statusInUse,
				LastStartedAt: time.Now().Add(-1000 * time.Second),
				Duration:      1800,
			},
			want: (1800 - 1000) * time.Second,
		},
		{
			m: &machine{
				Status:        statusIdle,
				LastStartedAt: time.Now().Add(-1000 * time.Second),
				Duration:      1800,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt.m.computeTimeLeft()
		if float64(tt.m.TimeLeft-tt.want) > 0.00001 {
			t.Errorf("got %v, want %v", tt.m.TimeLeft, tt.want)
		}
	}
}
