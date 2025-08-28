package fail

// ErrorUserMessage is an error type that provides a user-facing message.
//
// Implementations of this interface should return a concise, human-readable message
// suitable for display to end users. The message must not expose internal details
// or personally identifiable information (PII), and must not include details from
// any wrapped errors or underlying causes. It should represent only the primary error itself.
//
// Example usage:
//
//	type MyUserError struct{}
//	func (e *MyUserError) Error() string { return "internal: failed to connect to DB" }
//	func (e *MyUserError) ErrorUserMessage() string { return "Could not complete your request. Please try again." }
//
//	err := &MyUserError{}
//	msg := fail.UserMessage(err) // returns "Could not complete your request. Please try again."
type ErrorUserMessage interface {
	error

	// ErrorUserMessage returns the canonical user-facing message associated with this error.
	//
	// This message should be suitable for display to end users and must not
	// include the message of any wrapped error, error cause, internal details, or PII.
	ErrorUserMessage() string
}

// UserMessage returns the user-facing message for the provided error.
//
// This function determines the user message as follows:
//  1. If err is nil, it returns the empty string.
//  2. If err implements ErrorUserMessage, it returns the result of ErrorUserMessage().
//  3. Otherwise, it returns err.Error() (which may include internal details and is not guaranteed to be user-safe).
//
// This allows error types to specify custom user-facing messages, and for composed/multi-cause errors
// to propagate the most appropriate message for end users.
func UserMessage(err error) string {
	if err == nil {
		return ""
	}

	if message, ok := err.(ErrorUserMessage); ok {
		return message.ErrorUserMessage()
	}

	return err.Error()
}

// WithUserMessage returns a new error with the specified user-facing message attached.
//
// This function wraps an existing error with a user message string suitable for display to end users.
// If the provided error is nil, it returns nil. If the userMessage string is empty, the original error is returned unchanged.
//
// The returned error will implement the ErrorUserMessage interface, allowing retrieval of the user message via fail.UserMessage.
//
// Example:
//
//	err := fail.WithUserMessage(primaryErr, "Could not complete your request. Please try again.")
//
// The returned error will have the user message attached, which can be accessed using
// fail.UserMessage(err).
//
// Parameters:
//   - err:         The error to which the user message will be attached.
//   - userMessage: The user-facing message string to associate with the error.
//
// Returns:
//   - A new error with the user message attached, or nil if err is nil. If userMessage is empty, returns the original error.
func WithUserMessage(err error, userMessage string) error {
	if err == nil {
		return nil
	}

	if userMessage == "" {
		return err
	}

	return From(err).UserMsg(userMessage).asFail()
}
