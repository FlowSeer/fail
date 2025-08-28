package fail

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// ErrorTraceId is an error type that provides a trace ID associated with the error.
//
// Implementations of this interface should return a string representing the trace ID
// for distributed tracing purposes. The returned string may be empty if no trace ID is set.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorTraceId() string { return "abcdef1234567890" }
//
//	err := &MyError{}
//	traceId := fail.TraceId(err) // returns "abcdef1234567890"
type ErrorTraceId interface {
	error

	// ErrorTraceId returns the trace ID associated with this error.
	//
	// The returned string may be empty if no trace ID is set.
	ErrorTraceId() string
}

// TraceId returns the trace ID associated with the provided error, if any.
//
// This function attempts to extract the trace ID from the error as follows:
//  1. If err is nil, it returns an empty string.
//  2. If err implements ErrorTraceId, it returns the result of ErrorTraceId().
//  3. Otherwise, it returns an empty string.
//
// The returned string may be empty if no trace ID is set.
func TraceId(err error) string {
	if err == nil {
		return ""
	}

	if t, ok := err.(ErrorTraceId); ok {
		return t.ErrorTraceId()
	}

	return ""
}

// WithTraceId returns a new error with the specified trace ID attached.
//
// This function wraps an existing error with a trace ID string for distributed tracing.
// If the provided error is nil, it returns nil. If the trace ID string is empty, the original error is returned unchanged.
// If traceId is non-empty but not a valid hexadecimal trace.TraceID, the returned error will implement ErrorTraceId but return an empty trace ID.
//
// The resulting error will implement the ErrorTraceId interface, allowing retrieval of the trace ID via fail.TraceId.
//
// Example:
//
//	err := fail.WithTraceId(primaryErr, "abcdef1234567890")
//
// The returned error will have the trace ID attached, which can be accessed using
// fail.TraceId(err).
//
// Parameters:
//   - err:     The error to which the trace ID will be attached.
//   - traceId: The trace ID string to associate with the error.
//
// Returns:
//   - A new error with the trace ID attached, or nil if err is nil. If traceId is empty, returns the original error.
func WithTraceId(err error, traceId string) error {
	if err == nil {
		return nil
	}

	if traceId == "" {
		return err
	}

	return From(err).TraceId(traceId).asFail()
}

// TraceIdFromContext extracts the trace ID from the provided context using OpenTelemetry.
//
// This function returns the trace ID as a string from the current span in the context.
// If no span is present, the returned string will be empty.
//
// Example usage:
//
//	traceId := fail.TraceIdFromContext(ctx)
func TraceIdFromContext(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
