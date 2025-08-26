package fail

import (
	"time"

	"github.com/FlowSeer/wz/maps"
	"github.com/FlowSeer/wz/slices"
)

// Fail is a rich error type that implements all fail.* error interfaces.
//
// It provides support for error codes, exit codes, HTTP status codes, causes, associated errors,
// tags, and arbitrary attributes. This struct is intended to be used as the canonical error
// implementation for the fail package.
type Fail struct {
	time time.Time // Timestamp of when the error occurred

	msg     string // The main error message (required, never empty)
	userMsg string // Optional user-facing message

	domain         string // Domain of the error
	code           string // Application-specific error code
	exitCode       int    // Process exit code
	httpStatusCode int    // HTTP status code

	causes     []error // Direct causes of this error
	associated []error // Associated (but not causal) errors

	tags  map[string]struct{} // Set of string tags
	attrs map[string]any      // Arbitrary key-value attributes

	spanId  string // spanId is the unique identifier for the tracing span associated with this error.
	traceId string // traceId is the unique identifier for the tracing trace associated with this error.
}

// newFail creates a new Fail error with the given message.
//
// The message must not be an empty string. The returned Fail will have default values
// for code, exitCode, httpStatusCode, and empty tags/attributes.
func newFail(msg string) Fail {
	return Fail{
		msg:            msg,
		code:           DefaultErrorCode,
		exitCode:       DefaultExitCode,
		httpStatusCode: DefaultHttpStatusCode,
		tags:           make(map[string]struct{}),
		attrs:          make(map[string]any),
	}
}

// newFailFrom creates a new Fail error from an existing error.
//
// If the provided error is already a Fail, it returns a clone of that Fail.
// Otherwise, it constructs a new Fail by copying all available error details from the source error,
// including: message, user message, code, exit code, HTTP status code, causes, associated errors,
// tags, and attributes. Panics if err is nil.
func newFailFrom(err error) Fail {
	if err == nil {
		panic("cannot create a Fail from a nil error")
	}

	if f, ok := err.(Fail); ok {
		return f.Clone()
	}

	attrs := make(map[string]struct{})
	for _, t := range Tags(err) {
		attrs[t] = struct{}{}
	}

	return Fail{
		msg:            Message(err),
		userMsg:        UserMessage(err),
		code:           Code(err),
		exitCode:       ExitCode(err),
		httpStatusCode: HttpStatusCode(err),
		causes:         Causes(err),
		associated:     Associated(err),
		tags:           attrs,
		attrs:          Attributes(err),
	}
}

// Clone returns a deep copy of the Fail error.
//
// All fields are copied, including slices and maps, so that modifications to the
// returned Fail do not affect the original. This is useful for creating a new error
// instance based on an existing one, without sharing mutable state.
func (f Fail) Clone() Fail {
	return Fail{
		msg:            f.msg,
		userMsg:        f.userMsg,
		code:           f.code,
		exitCode:       f.exitCode,
		httpStatusCode: f.httpStatusCode,
		causes:         slices.Clone(f.causes),
		associated:     slices.Clone(f.associated),
		tags:           maps.Clone(f.tags),
		attrs:          maps.Clone(f.attrs),
	}
}

// Error returns the main error message.
func (f Fail) Error() string {
	return f.msg
}

// ErrorCauses returns the direct causes of this error.
//
// Implements ErrorCauses interface.
func (f Fail) ErrorCauses() []error {
	return f.causes
}

// ErrorAssociated returns the associated (non-causal) errors.
//
// Implements ErrorAssociated interface. The returned slice is a copy.
func (f Fail) ErrorAssociated() []error {
	return slices.Clone(f.associated)
}

// ErrorCode returns the application-specific error code.
//
// Implements ErrorCode interface.
func (f Fail) ErrorCode() string {
	return f.code
}

// ErrorExitCode returns the process exit code for this error.
//
// Implements ErrorExitCode interface.
func (f Fail) ErrorExitCode() int {
	return f.exitCode
}

// ErrorHttpStatusCode returns the HTTP status code for this error.
//
// Implements ErrorHttpStatusCode interface.
func (f Fail) ErrorHttpStatusCode() int {
	return f.httpStatusCode
}

// ErrorMessage returns the main error message.
//
// Implements ErrorMessage interface.
func (f Fail) ErrorMessage() string {
	return f.msg
}

// ErrorUserMessage returns the user-facing error message, if any.
//
// Implements ErrorUserMessage interface.
func (f Fail) ErrorUserMessage() string {
	return f.userMsg
}

// ErrorTags returns a slice of tags associated with this error.
//
// Implements ErrorTags interface. The returned slice is a copy.
func (f Fail) ErrorTags() []string {
	return slices.Collect(maps.Keys(f.tags))
}

// ErrorAttributes returns a copy of the attributes map for this error.
//
// Implements ErrorAttributes interface.
func (f Fail) ErrorAttributes() map[string]any {
	return maps.Clone(f.attrs)
}

// ErrorTime returns the timestamp of when the error occurred.
//
// Implements ErrorTime interface.
func (f Fail) ErrorTime() time.Time {
	return f.time
}
