package fail

import "time"

// ErrorTime is an interface for errors that have an associated time.
//
// Implementing this interface allows an error to expose a timestamp
// (such as when the error occurred, was logged, or was created).
// This can be useful for diagnostics, logging, or error reporting.
//
// Example:
//
//	type MyError struct {
//	    when time.Time
//	}
//	func (e *MyError) Error() string { return "something happened" }
//	func (e *MyError) ErrorTime() time.Time { return e.when }
type ErrorTime interface {
	// ErrorTime returns the time associated with this error.
	ErrorTime() time.Time
}

// Time returns the time associated with the provided error, if available.
//
// If the error implements the ErrorTime interface, its ErrorTime() value is returned.
// If err is nil or does not implement ErrorTime, the zero time (time.Time{}) is returned.
//
// This function is useful for retrieving timestamps from errors that carry time information,
// such as when the error occurred or was recorded.
func Time(err error) time.Time {
	if err == nil {
		return time.Time{}
	}

	if t, ok := err.(ErrorTime); ok {
		return t.ErrorTime()
	}

	return time.Time{}
}

// WithTime returns a new error with the specified time.Time value attached.
//
// This function wraps an existing error with a timestamp, which can be useful for
// diagnostics, logging, or error reporting. If the provided error is nil, it returns nil.
// If the provided time is the zero value (t.IsZero()), the original error is returned unchanged.
//
// The returned error will implement the ErrorTime interface, allowing retrieval of the
// associated time via fail.Time(err).
//
// Example:
//
//	err := fail.WithTime(primaryErr, time.Now())
//
// The returned error will have the time attached, which can be accessed using
// fail.Time(err).
//
// Parameters:
//   - err: The error to which the time will be attached.
//   - t:   The time.Time value to associate with the error.
//
// Returns:
//   - A new error with the time attached, or nil if err is nil. If t is zero, returns the original error.
func WithTime(err error, t time.Time) error {
	if err == nil {
		return nil
	}

	if t.IsZero() {
		return err
	}

	return From(err).Time(t).asFail()
}

// WithTimeNow returns a new error with the current time attached.
//
// This is a convenience function equivalent to calling WithTime(err, time.Now()).
// If the provided error is nil, it returns nil.
//
// Example:
//
//	err := fail.WithTimeNow(primaryErr)
//
// The returned error will have the current time attached, which can be accessed using
// fail.Time(err).
//
// Parameters:
//   - err: The error to which the current time will be attached.
//
// Returns:
//   - A new error with the current time attached, or nil if err is nil.
func WithTimeNow(err error) error {
	return WithTime(err, time.Now())
}
