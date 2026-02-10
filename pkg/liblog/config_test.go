package liblog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		c                       Config
		expectedErrorOccurrence bool
	}{
		{
			name: "test case (empty Formatter config specified)",
			c: Config{
				Level:   InfoLevel,
				Outputs: []Output{StdOut},
			},
			expectedErrorOccurrence: true,
		},
		{
			name: "test case (valid config specified)",
			c: Config{
				Level:     InfoLevel,
				Outputs:   []Output{StdOut},
				Formatter: FormatJSON,
			},
			expectedErrorOccurrence: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				actualError           = tt.c.Validate()
				actualErrorOccurrence = actualError != nil
			)

			assert.Equal(
				t,
				tt.expectedErrorOccurrence,
				actualErrorOccurrence,
				fmt.Sprintf("actualError is [%v]", actualError),
			)
		})
	}
}

func TestLevel_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		l             Level
		expectedError error
	}{
		{
			name:          "test-1 (empty log level specified)",
			l:             "",
			expectedError: errEmptyLogLevel,
		},
		{
			name:          "test-2 (non existed log level specified)",
			l:             "NonExistedLevel", // this type doesn't exist, so we should get an error
			expectedError: errInvalidLogLevel,
		},
		{
			name:          "test-3 (valid log level specified)",
			l:             InfoLevel,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualError := tt.l.Validate()

			assert.Equal(
				t,
				tt.expectedError,
				actualError,
			)
		})
	}
}

func TestOutput_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		o             Output
		expectedError error
	}{
		{
			name:          "test-1 (empty log output specified)",
			o:             "",
			expectedError: errEmptyLoggerOutput,
		},
		{
			name:          "test-2 (non existed log output specified)",
			o:             "NonExistedOutput",
			expectedError: nil,
		},
		{
			name:          "test-3 (valid log output specified)",
			o:             StdOut,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualError := tt.o.Validate()

			assert.Equal(
				t,
				tt.expectedError,
				actualError,
			)
		})
	}
}

func TestOutputs_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		o             Outputs
		expectedError error
	}{
		{
			name:          "test-1 (empty log output slice specified)",
			o:             Outputs{},
			expectedError: errEmptyLoggerOutput,
		},
		{
			name:          "test-2 (duplicate log output specified)",
			o:             Outputs{StdOut, StdOut},
			expectedError: errDuplicateLogOutput,
		},
		{
			name:          "test-3 (valid log output slice specified)",
			o:             Outputs{StdOut},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualError := tt.o.Validate()

			assert.Equal(
				t,
				tt.expectedError,
				actualError,
			)
		})
	}
}

func TestFormat_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		f             Format
		expectedError error
	}{
		{
			name:          "test-1 (empty log format specified)",
			f:             "",
			expectedError: errEmptyLogFormat,
		},
		{
			name:          "test-2 (non existed log format specified)",
			f:             "NonExistedFormat",
			expectedError: errInvalidLogFormat,
		},
		{
			name:          "test-3 (valid log format specified)",
			f:             FormatJSON,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualError := tt.f.Validate()

			assert.Equal(
				t,
				tt.expectedError,
				actualError,
			)
		})
	}
}
