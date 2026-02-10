package liblog

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	t.Parallel()

	mockLogger := NewMockLogger(t)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		expectedLogger Logger
		expectedExists bool
	}{
		{
			name: "test-1 (ctx without logger specified)",
			args: args{
				ctx: context.Background(),
			},
			expectedLogger: nil,
			expectedExists: false,
		},
		{
			name: "test-2 (ctx with logger specified)",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey{}, mockLogger),
			},
			expectedLogger: mockLogger,
			expectedExists: true,
		},
	}
	for _, tt := range tests {
		pp := tt
		t.Run(pp.name, func(t *testing.T) {
			t.Parallel()

			actualLogger, actualExists := From(pp.args.ctx)
			assert.Equal(
				t,
				pp.expectedLogger,
				actualLogger,
			)
			assert.Equal(
				t,
				pp.expectedExists,
				actualExists,
			)
		})
	}
}

func TestMust(t *testing.T) {
	t.Parallel()

	mockLogger := NewMockLogger(t)

	type args struct {
		ctx context.Context
	}

	panicTestNameForRecovery := "test-1 (ctx without logger specified)"

	tests := []struct {
		name           string
		args           args
		expectedLogger Logger
	}{
		{
			name: panicTestNameForRecovery,
			args: args{
				ctx: context.Background(),
			},
			expectedLogger: nil, // should panic
		},
		{
			name: "test-2 (ctx with logger specified)",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey{}, mockLogger),
			},
			expectedLogger: mockLogger,
		},
	}
	for _, tt := range tests {
		pp := tt
		t.Run(pp.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); r != nil {
					assert.Equal(
						t,
						panicTestNameForRecovery,
						pp.name,
						fmt.Sprintf("should panic for %s", pp.name),
					)

					return
				}
			}()

			actualLogger := Must(pp.args.ctx)
			assert.Equal(
				t,
				pp.expectedLogger,
				actualLogger,
			)
		})
	}
}

func TestWith(t *testing.T) {
	t.Parallel()

	mockLogger := NewMockLogger(t)

	type args struct {
		ctx    context.Context
		logger Logger
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-1 (ctx and logger specified)",
			args: args{
				ctx:    context.Background(),
				logger: mockLogger,
			},
		},
	}
	for _, tt := range tests {
		pp := tt
		t.Run(pp.name, func(t *testing.T) {
			t.Parallel()

			actualCtx := With(pp.args.ctx, pp.args.logger)

			actualLogger, exists := From(actualCtx)

			assert.Equal(
				t,
				pp.args.logger,
				actualLogger,
			)

			assert.Equal(
				t,
				true,
				exists,
			)
		})
	}
}
