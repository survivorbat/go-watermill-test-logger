package watertestlogger

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	testingi "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
)

func TestTestLogAdaptor_RespectsLogLevel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		logLevel      watermill.LogLevel
		expectedCalls [][]any
	}{
		"trace": {
			logLevel: watermill.TraceLogLevel,
			expectedCalls: [][]any{
				{"trace", "map[foo:bar]"},
				{"debug", "map[foo:bar]"},
				{"info", "map[foo:bar]"},
				{"error", assert.AnError.Error(), "map[foo:bar]"},
			},
		},
		"debug": {
			logLevel: watermill.DebugLogLevel,
			expectedCalls: [][]any{
				{"debug", "map[foo:bar]"},
				{"info", "map[foo:bar]"},
				{"error", assert.AnError.Error(), "map[foo:bar]"},
			},
		},
		"info": {
			logLevel: watermill.InfoLogLevel,
			expectedCalls: [][]any{
				{"info", "map[foo:bar]"},
				{"error", assert.AnError.Error(), "map[foo:bar]"},
			},
		},
		"error": {
			logLevel: watermill.ErrorLogLevel,
			expectedCalls: [][]any{
				{"error", assert.AnError.Error(), "map[foo:bar]"},
			},
		},
	}

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			testT := new(testT)
			adaptor := NewTestAdaptor(testT, false, testData.logLevel)

			fields := watermill.LogFields{"foo": "bar"}

			// Act
			adaptor.Trace("trace", fields)
			adaptor.Debug("debug", fields)
			adaptor.Info("info", fields)
			adaptor.Error("error", assert.AnError, fields)

			// Assert
			assert.Equal(t, testData.expectedCalls, testT.LogCalls)
		})
	}
}

func TestTestLogAdaptor_WritesErrorsIfConfigured(t *testing.T) {
	t.Parallel()
	// Arrange
	testT := new(testT)
	adaptor := NewTestAdaptor(testT, true, watermill.ErrorLogLevel)

	fields := watermill.LogFields{"foo": "bar"}

	// Act
	adaptor.Error("error", assert.AnError, fields)

	// Assert
	expectedCalls := [][]any{
		{"error", assert.AnError.Error(), "map[foo:bar]"},
	}

	assert.Equal(t, expectedCalls, testT.ErrorCalls)
}

func TestTestLogAdaptor_With_SavesFields(t *testing.T) {
	t.Parallel()
	// Arrange
	testT := new(testT)
	adaptor := NewTestAdaptor(testT, true, watermill.InfoLogLevel)

	fieldsA := watermill.LogFields{"foo": "bar"}
	fieldsB := watermill.LogFields{"baz": "bar"}

	// Act
	adaptor.With(fieldsA).With(fieldsB)

	// Assert
	expectedFields := watermill.LogFields{
		"foo": "bar",
		"baz": "bar",
	}
	assert.Equal(t, expectedFields, adaptor.withFields)
}

// Mocks

var _ testingi.T = new(testT)

type testT struct {
	testing.T

	LogCalls   [][]any
	ErrorCalls [][]any
}

func (t *testT) Log(args ...any) {
	t.LogCalls = append(t.LogCalls, args)
}

func (t *testT) Error(args ...any) {
	t.ErrorCalls = append(t.ErrorCalls, args)
}
