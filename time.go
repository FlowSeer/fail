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
