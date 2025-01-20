package watertestlogger

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	testingi "github.com/mitchellh/go-testing-interface"
)

// Compile-time interface checks
var _ watermill.LoggerAdapter = new(TestLogAdaptor)

// NewTestAdaptor instantiates a new adaptor that writes logs to the given test instance. If
// useError is set to true, it will also cause the test to fail.
func NewTestAdaptor(t testingi.T, useError bool, logLevel watermill.LogLevel) *TestLogAdaptor {
	return &TestLogAdaptor{
		UseError:   useError,
		LogLevel:   logLevel,
		t:          t,
		withFields: make(watermill.LogFields),
	}
}

// TestLogAdaptor implements watermill.LoggerAdapter and can be used to capture logs
// in test code.
type TestLogAdaptor struct {
	// UseError will instruct the adaptor to use t.Error when logging fields, causing
	// tests to fail if any error were logged. Exposed to allow users to change the value
	// during a test.
	UseError bool

	// LogLevel is uses to determine what we should log and what we should skip. Exposed to allow users to change the value
	// during a test.
	LogLevel watermill.LogLevel

	// t is the test instance that we'll write logs to.
	t testingi.T

	// withFields are saved from the With call and appended later
	withFields watermill.LogFields
}

func (t *TestLogAdaptor) log(level watermill.LogLevel, msg string, fields watermill.LogFields) {
	if t.LogLevel > level {
		return
	}

	if fields == nil {
		fields = make(watermill.LogFields)
	}

	mergedFields := t.withFields
	for key, field := range fields {
		mergedFields[key] = field
	}

	t.t.Log(msg, fmt.Sprintf("%+v", mergedFields))
}

// Debug implements watermill.LoggerAdapter.
func (t *TestLogAdaptor) Debug(msg string, fields watermill.LogFields) {
	t.log(watermill.DebugLogLevel, msg, fields)
}

// Info implements watermill.LoggerAdapter.
func (t *TestLogAdaptor) Info(msg string, fields watermill.LogFields) {
	t.log(watermill.InfoLogLevel, msg, fields)
}

// Trace implements watermill.LoggerAdapter.
func (t *TestLogAdaptor) Trace(msg string, fields watermill.LogFields) {
	t.log(watermill.TraceLogLevel, msg, fields)
}

// Error implements watermill.LoggerAdapter.
func (t *TestLogAdaptor) Error(msg string, err error, fields watermill.LogFields) {
	if fields == nil {
		fields = make(watermill.LogFields)
	}

	mergedFields := t.withFields
	for key, field := range fields {
		mergedFields[key] = field
	}

	if t.UseError {
		t.t.Error(msg, err.Error(), fmt.Sprintf("%+v", mergedFields))
		return
	}

	t.t.Log(msg, err.Error(), fmt.Sprintf("%+v", mergedFields))
}

// With implements watermill.LoggerAdapter.
func (t *TestLogAdaptor) With(fields watermill.LogFields) watermill.LoggerAdapter {
	mergedFields := t.withFields
	for key, field := range fields {
		mergedFields[key] = field
	}

	return &TestLogAdaptor{
		UseError:   t.UseError,
		LogLevel:   t.LogLevel,
		t:          t.t,
		withFields: mergedFields,
	}
}
