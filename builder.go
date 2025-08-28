package fail

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Builder is a fluent builder for constructing Fail errors with additional context,
// such as tags, attributes, causes, associated errors, codes, user messages, and tracing info.
//
// The Builder implements a fluent interface that allows chaining method calls to construct
// rich error objects. Each method returns the Builder itself, enabling method chaining.
// The Msg() and Msgf() methods are terminal and return the complete Fail error that
// implements all the fail.* error interfaces.
//
// Example usage:
//
//	err := fail.New().
//		UserMsg("Unable to process your request. Please try again.").
//		Code("DB_CONNECTION_ERROR").
//		Domain(fail.DomainDatabase).
//		ExitCode(1).
//		HttpStatusCode(503).
//		Tag(fail.TagNetwork, fail.TagTimeout, fail.TagDatabase).
//		Attribute("host", "db.example.com").
//		Attribute("port", 5432).
//		Cause(io.EOF).
//		Associate(loggingError).
//		TraceId("abcdef1234567890").
//		SpanId("1234567890abcdef").
//		Msg("database connection failed")
type Builder Fail

// New creates a new Builder with an empty message.
//
// The returned Builder will have default values for code (DefaultErrorCode),
// exit code (DefaultExitCode), and HTTP status code (DefaultHttpStatusCode).
// The message must be set using Msg() or Msgf() to complete the error construction.
// If no message is set, the message will be set to fail.EmptyMessage.
// The timestamp will be automatically set to the current time when the error is built
// if not explicitly set using the Time() method.
//
// Example:
//
//	builder := fail.New()
func New() Builder {
	return Builder(newFail(""))
}

// From creates a new Builder initialized from an existing error.
//
// If the provided error is already a Fail, it returns a new Builder populated with the same details.
// Otherwise, it constructs a new Builder by extracting all available error details from the source error,
// including: message, user message, code, exit code, HTTP status code, causes, associated errors,
// tags, and attributes. Panics if err is nil.
//
// Example:
//
//	err := someFunction()
//	failErr := fail.From(err).Msg("operation failed")
func From(err error) Builder {
	if err == nil {
		panic("cannot create a Fail from a nil error")
	}

	if f, ok := err.(Fail); ok {
		return Builder(f.Clone())
	}

	attrs := make(map[string]struct{})
	for _, t := range Tags(err) {
		attrs[t] = struct{}{}
	}

	return Builder(Fail{
		msg:            Message(err),
		userMsg:        UserMessage(err),
		code:           Code(err),
		exitCode:       ExitCode(err),
		httpStatusCode: HttpStatusCode(err),
		causes:         Causes(err),
		associated:     Associated(err),
		tags:           attrs,
		attrs:          Attributes(err),
	})
}

// Time sets the timestamp for when the error occurred.
//
// If the provided time is not the zero value and is not in the future, it will be set as the error's time.
// If no time is set or if the set time is in the future, the timestamp will be automatically set to the current time when the error is built using Msg() or Msgf().
//
// Example:
//
//	err := fail.New().
//		Time(time.Now()).
//		Msg("operation failed")
func (b Builder) Time(t time.Time) Builder {
	if !t.IsZero() && time.Now().After(t) {
		b.time = t
	}

	return b
}

// Associate adds one or more associated errors to the builder.
// Associated errors are related errors that provide additional context but are not direct causes.
//
// Associated errors implement the ErrorAssociated interface and represent errors that are
// related to the current error but are not direct causes. Examples include:
//   - Errors that occur during cleanup after a primary error
//   - Logging failures during error reporting
//   - Multiple independent errors in batch operations
//   - Errors from different nodes in distributed systems
//
// The associated errors will be accessible via the Associated() function on the built error.
//
// Example:
//
//	err := fail.New().
//		Cause(networkError).
//		Associate(diskWriteError, loggingError).
//		Msg("file upload failed")
func (b Builder) Associate(errs ...error) Builder {
	return b.AssociateSlice(errs)
}

// AssociateSlice adds a slice of associated errors to the builder.
//
// An associated error is a related error that may be encountered while handling the cause,
// or an error that is not the direct cause of the error but is still relevant to the error.
// This method is useful when you already have a slice of associated errors.
//
// The associated errors will be accessible via the Associated() function on the built error.
//
// Example:
//
//	associatedErrors := []error{diskWriteError, loggingError}
//	err := fail.New().
//		AssociateSlice(associatedErrors).
//		Msg("file upload failed")
func (b Builder) AssociateSlice(errs []error) Builder {
	for _, err := range errs {
		if err != nil {
			b.associated = append(b.associated, err)
		}
	}
	return b
}

// Cause adds one or more cause errors to the builder.
//
// A cause error is an error that directly led to this error and represent the underlying reasons for the current error.
// The Causes() function will return these errors when called on the built error.
//
// Cause errors are different from associated errors in that they represent the direct
// chain of causality, while associated errors are related but not causal.
//
// Example:
//
//	err := fail.New().
//		Cause(connectionError, queryError).
//		Msg("database operation failed")
func (b Builder) Cause(errs ...error) Builder {
	return b.CauseSlice(errs)
}

// CauseSlice adds a slice of cause errors to the builder.
//
// A cause error is an error that directly led to this error.
//
// Example:
//
//	causeErrors := []error{connectionError, queryError}
//	err := fail.New().
//		CauseSlice(causeErrors).
//		Msg("database operation failed")
func (b Builder) CauseSlice(errs []error) Builder {
	for _, err := range errs {
		if err != nil {
			b.causes = append(b.causes, err)
		}
	}
	return b
}

// Tag adds one or more tags to the builder.
//
// A tag is a string label that can be used for categorization or filtering and provide a way to categorize errors for logging, monitoring, or error handling purposes.
// Common tags include domain names, error types, or system components.
//
// Example:
//
//	err := fail.New().
//		Tag("api", "network", "timeout").
//		Msg("API request failed")
func (b Builder) Tag(tags ...string) Builder {
	return b.TagSlice(tags)
}

// TagSlice adds a slice of tags to the builder.
//
// A tag is a string label that can be used for categorization or filtering.
//
// Example:
//
//	tags := []string{"database", "connection", "timeout"}
//	err := fail.New().
//		TagSlice(tags).
//		Msg("database connection failed")
func (b Builder) TagSlice(tags []string) Builder {
	for _, tag := range tags {
		if tag != "" {
			b.tags[tag] = struct{}{}
		}
	}
	return b
}

// Domain sets the domain for the error being built.
//
// The domain is a string that categorizes the error by its source or type, such as "network", "database", or "validation".
// Domains are useful for grouping, filtering, and handling errors in a structured way throughout your application.
// If the provided domain is an empty string, the builder's domain is not changed.
//
// Example:
//
//	err := fail.New().
//		Domain(fail.DomainDatabase).
//		Msg("failed to connect to database")
func (b Builder) Domain(domain string) Builder {
	if domain != "" {
		b.domain = domain
	}

	return b
}

// Attribute adds a key-value attribute to the builder.
//
// An attribute is a key-value pair that provides additional structured context and allow you to attach arbitrary data to errors for debugging, logging, or monitoring purposes.
//
// Attributes can contain any type of value (interface{}), making them flexible for storing various types of contextual information such as request IDs, user IDs, timestamps, or other relevant data.
//
// Example:
//
//	err := fail.New().
//		Attribute("user_id", "12345").
//		Attribute("request_id", "req-abc-123").
//		Attribute("attempt_count", 3).
//		Msg("user authentication failed")
func (b Builder) Attribute(key string, value any) Builder {
	return b.AttributeMap(map[string]any{key: value})
}

// AttributeMap adds a map of key-value attributes to the builder.
//
// An attribute is a key-value pair that provides additional structured context.
//
// Example:
//
//	attrs := map[string]any{
//		"user_id": "12345",
//		"request_id": "req-abc-123",
//		"attempt_count": 3,
//	}
//	err := fail.New().
//		AttributeMap(attrs).
//		Msg("user authentication failed")
func (b Builder) AttributeMap(attrs map[string]any) Builder {
	for key, value := range attrs {
		if key != "" && value != nil {
			b.attrs[key] = value
		}
	}
	return b
}

// Code sets a string code for the error, such as an error type or identifier.
//
// A code is a string that can be used to identify the error and should be a stable, concise string that uniquely identifies the type or category of the error.
// The code must not contain whitespace or special characters—only letters, numbers, and underscores are allowed.
//
// Example:
//
//	err := fail.New().
//		Code("VALIDATION_ERROR").
//		Msg("invalid input provided")
func (b Builder) Code(code string) Builder {
	if code != "" {
		b.code = code
	}
	return b
}

// ExitCode sets a process exit code for the error, if greater than zero.
//
// The exit code represents the process exit status that should be used when this error occurs.
// Only positive values are accepted; negative or zero values are ignored.
//
// Example:
//
//	err := fail.New().
//		ExitCode(2).
//		Msg("configuration file not found")
func (b Builder) ExitCode(exitCode int) Builder {
	if exitCode > 0 {
		b.exitCode = exitCode
	}
	return b
}

// HttpStatusCode sets an HTTP status code for the error, if in the 400-599 range.
//
// The HTTP status code represents the HTTP response status that should be returned when this error occurs in an HTTP context.
// Only status codes in the 400-599 range (client and server errors) are accepted.
//
// Example:
//
//	err := fail.New().
//		HttpStatusCode(404).
//		Msg("user not found")
func (b Builder) HttpStatusCode(httpStatusCode int) Builder {
	if httpStatusCode >= 400 && httpStatusCode < 600 {
		b.httpStatusCode = httpStatusCode
	}
	return b
}

// TraceId sets the trace ID for distributed tracing, if the string is a valid hex trace ID.
//
// The trace ID is used for distributed tracing to correlate errors across different services and components.
// The trace ID must be a valid hexadecimal string representation of a trace ID.
//
// Example:
//
//	err := fail.New().
//		TraceId("abcdef1234567890abcdef1234567890").
//		Msg("request processing failed")
func (b Builder) TraceId(traceId string) Builder {
	t, err := trace.TraceIDFromHex(traceId)
	if err == nil {
		b.traceId = t.String()
	}
	return b
}

// SpanId sets the span ID for distributed tracing, if the string is a valid hex span ID.
//
// The span ID is used for distributed tracing to identify specific spans within a trace.
// The span ID must be a valid hexadecimal string representation of a span ID.
//
// Example:
//
//	err := fail.New().
//		SpanId("1234567890abcdef").
//		Msg("database query failed")
func (b Builder) SpanId(spanId string) Builder {
	s, err := trace.SpanIDFromHex(spanId)
	if err == nil {
		b.spanId = s.String()
	}
	return b
}

// Context extracts tags, attributes, span ID, and trace ID from the provided context.Context and adds them to the builder, if present.
//
// This method automatically extracts error-related information from the context using the following functions:
//   - TagsFromContext(): Extracts tags stored in the context
//   - AttributesFromContext(): Extracts attributes stored in the context
//   - SpanIdFromContext(): Extracts the span ID from OpenTelemetry span in the context
//   - TraceIdFromContext(): Extracts the trace ID from OpenTelemetry span in the context
//
// This is useful for propagating error context through request lifecycles or operation
// chains without manually passing each component.
//
// Example:
//
//	ctx := context.Background()
//	ctx = fail.ContextWithTags(ctx, []string{"api", "v1"})
//	ctx = fail.ContextWithAttributes(ctx, map[string]any{"user_id": "123"})
//
//	err := fail.New().
//		Context(ctx).
//		Msg("request failed")
func (b Builder) Context(ctx context.Context) Builder {
	res := b

	tags := TagsFromContext(ctx)
	if tags != nil {
		res = res.TagSlice(tags)
	}

	attrs := AttributesFromContext(ctx)
	if len(attrs) > 0 {
		res = res.AttributeMap(attrs)
	}

	spanId := SpanIdFromContext(ctx)
	if spanId != "" {
		res = res.SpanId(spanId)
	}

	traceId := TraceIdFromContext(ctx)
	if traceId != "" {
		res = res.TraceId(traceId)
	}

	return res
}

// UserMsg sets a user-facing message for the error.
//
// The user message implements the ErrorUserMessage interface and provides a concise,
// human-readable message suitable for display to end users. The message must not expose
// internal details, personally identifiable information (PII), or details from wrapped
// errors or underlying causes.
//
// The UserMessage() function will return this message when called on the built error.
//
// Example:
//
//	err := fail.New().
//		UserMsg("We're experiencing technical difficulties. Please try again later.").
//		Msg("database connection failed: connection refused")
func (b Builder) UserMsg(userMsg string) Builder {
	if userMsg != "" {
		b.userMsg = userMsg
	}
	return b
}

// UserMsgf sets a formatted user-facing message for the error.
//
// This method is a convenience wrapper around UserMsg() that allows formatting the
// user message using fmt.Sprintf. The formatted message must still adhere to the
// requirements for user messages (no internal details, no PII, etc.).
//
// Example:
//
//	err := fail.New().
//		UserMsgf("Too many requests. Please wait %d seconds before trying again.", 60).
//		Msg("rate limit exceeded")
func (b Builder) UserMsgf(format string, args ...any) Builder {
	return b.UserMsg(fmt.Sprintf(format, args...))
}

// Msg sets a developer-facing message for the error and returns the complete Fail error.
//
// The developer message is the main error message and is required.
// If omitted, the message will be set to fail.EmptyMessage.
// It provides a concise, stable, and programmatically useful message that describes only the primary error itself, without including details from any wrapped errors or underlying causes.
// This method is terminal and completes the error construction.
//
// Example:
//
//	err := fail.New().
//		Code("DB_CONNECTION_ERROR").
//		Msg("database connection failed")
func (b Builder) Msg(msg string) error {
	if msg != "" {
		b.msg = msg
	} else {
		b.msg = EmptyMessage
	}

	if b.time.IsZero() || b.time.After(time.Now()) {
		b.time = time.Now()
	}

	return Fail(b)
}

// Msgf sets a formatted developer-facing message for the error and returns the complete Fail error.
//
// This method is a convenience wrapper around Msg() that allows formatting the
// developer message using fmt.Sprintf. The formatted message should still be concise
// and stable, suitable for logs or diagnostics.
// This method is terminal and completes the error construction.
//
// Example:
//
//	err := fail.New().
//		Code("DB_CONNECTION_ERROR").
//		Msgf("failed to connect to database %s on port %d", "localhost", 5432)
func (b Builder) Msgf(format string, args ...any) error {
	return b.Msg(fmt.Sprintf(format, args...))
}

// asFail returns the Builder as a Fail error as-is
func (b Builder) asFail() Fail {
	return Fail(b)
}
