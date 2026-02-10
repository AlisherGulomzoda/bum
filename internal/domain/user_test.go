package domain

import (
	"testing"
)

//nolint:nolintlint,all // it's ok
func TestUser_ValidateGender(t *testing.T) {
	type args struct {
		gender string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "male", args: struct{ gender string }{gender: "male"}, want: true},
		{name: "female", args: struct{ gender string }{gender: "female"}, want: true},
		{name: "unknown gender", args: struct{ gender string }{gender: "xyz"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Gender(tt.args.gender).Validate(); got != tt.want {
				t.Errorf("ValidateGender() = %v, want %v", got, tt.want)
			}
		})
	}
}
