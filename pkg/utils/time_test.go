package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	str := "2023-10-15T00:45:49+05:00"

	now, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Errorf("failed to parse time: %v", err)
		return
	}

	tests := []struct {
		name    string
		t       RFC3339Time
		want    []byte
		wantErr bool
	}{
		{
			"test1",
			RFC3339Time(now),
			[]byte(fmt.Sprintf("%q", str)),
			false,
		},
	}
	for _, tt := range tests {
		pp := tt
		t.Run(pp.name, func(t *testing.T) {
			t.Parallel()

			got, err := pp.t.MarshalJSON()
			if (err != nil) != pp.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, pp.wantErr)
				return
			}

			if !reflect.DeepEqual(got, pp.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(pp.want))
			}
		})
	}
}
