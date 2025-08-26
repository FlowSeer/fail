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
