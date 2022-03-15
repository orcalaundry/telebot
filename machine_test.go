package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_MachineUnmarshal(t *testing.T) {
	tests := []struct {
		json string   // the json string to unmarshal
		want *machine // machine value expected
		// lastStartedAtRFC3339 is a formatted string of the machine.LastStartedAt
		// time field, which helps get around some inconveniences here.
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
				Name:     "Coin washer",
				Floor:    14,
				Type:     washer,
				Status:   statusIdle,
				TimeLeft: 0,
			},
			lastStartedAtRFC3339: "1970-01-01T00:00:00+00:00",
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
			if !(reflect.DeepEqual(m, *tt.want)) {
				t.Errorf("got %+v, want %+v", m, tt.want)
			}
		})
	}
}
