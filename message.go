package fail

// ErrorMessage is an error type that provides a canonical, programmatic error message,
// distinct from the standard Error() string.
//
// The key difference between ErrorMessage() and Error() is that ErrorMessage() must return
// a concise, stable, and programmatically useful message that describes only the primary
// error itself, without including details from any wrapped errors or underlying causes.
// In contrast, Error() may include additional context, formatting, or information from
// wrapped errors, and is generally intended for human consumption (e.g., logs, error chains).
//
// ErrorMessage() is intended for logs, diagnostics, or programmatic error handling where
// a stable identifier for the error is needed, and is not meant for direct display to end users.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong: " + e.detail }
//	func (e *MyError) ErrorMessage() string { return "something went wrong" }
type ErrorMessage interface {
	error

	// ErrorMessage returns the canonical error message associated with this error.
	// This message should be suitable for logs or diagnostics and must not
	// include the message of any wrapped error or error cause.
	// Unlike Error(), this should be stable and not include context from other errors.
	ErrorMessage() string
}

// Message returns the canonical programmatic message for the provided error.
//
// This message is not intended to be user-facing, but rather for logs or diagnostics.
// If the error implements ErrorMessage, its ErrorMessage() is returned.
// Otherwise, err.Error() is returned. Note that Error() may include additional context
// or information from wrapped errors, while ErrorMessage() is intended to be stable and minimal.
func Message(err error) string {
	if err == nil {
		return ""
	}

	if message, ok := err.(ErrorMessage); ok {
		return message.ErrorMessage()
	}

	return err.Error()
}
