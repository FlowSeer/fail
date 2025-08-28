package fail

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// ErrorSpanId is an error type that provides a span ID associated with the error.
//
// Implementations of this interface should return a string representing the span ID
// for distributed tracing purposes. The returned string may be empty if no span ID is set.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorSpanId() string { return "1234567890abcdef" }
//
//	err := &MyError{}
//	spanId := fail.SpanId(err) // returns "1234567890abcdef"
type ErrorSpanId interface {
	error

	// ErrorSpanId returns the span ID associated with this error.
	//
	// The returned string may be empty if no span ID is set.
	ErrorSpanId() string
}

// SpanId returns the span ID associated with the provided error, if any.
//
// This function attempts to extract the span ID from the error as follows:
//  1. If err is nil, it returns an empty string.
//  2. If err implements ErrorSpanId, it returns the result of ErrorSpanId().
//  3. Otherwise, it returns an empty string.
//
// The returned string may be empty if no span ID is set.
func SpanId(err error) string {
	if err == nil {
		return ""
	}

	if span, ok := err.(ErrorSpanId); ok {
		return span.ErrorSpanId()
	}

	return ""
}

// WithSpanId returns a new error with the specified span ID attached.
//
// This function wraps an existing error with a span ID string for distributed tracing.
// If the provided error is nil, it returns nil. If the span ID string is empty, the original error is returned unchanged.
// If spanId is non-empty but not a valid hexadecimal trace.SpanID, the returned error will implement ErrorSpanId but return an empty span ID.
//
// The resulting error will implement the ErrorSpanId interface, allowing retrieval of the span ID via fail.SpanId.
//
// Example:
//
//	err := fail.WithSpanId(primaryErr, "1234567890abcdef")
//
// The returned error will have the span ID attached, which can be accessed using
// fail.SpanId(err).
//
// Parameters:
//   - err:    The error to which the span ID will be attached.
//   - spanId: The span ID string to associate with the error.
//
// Returns:
//   - A new error with the span ID attached, or nil if err is nil. If spanId is empty, returns the original error.
func WithSpanId(err error, spanId string) error {
	if err == nil {
		return nil
	}

	if spanId == "" {
		return err
	}

	return From(err).SpanId(spanId).asFail()
}

// SpanIdFromContext extracts the span ID from the provided context using OpenTelemetry.
//
// This function returns the span ID as a string from the current span in the context.
// If no span is present, the returned string will be empty.
//
// Example usage:
//
//	spanId := fail.SpanIdFromContext(ctx)
func SpanIdFromContext(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().SpanID().String()
}
