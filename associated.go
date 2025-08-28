package fail

// ErrorAssociated is an error type that provides a list of associated errors.
//
// Associated errors are errors that are related to the current error, but are not
// direct causes. For example, in a batch operation where multiple errors occur
// independently, each error may be associated with the others, but not caused by them.
// Another example is when an error occurs while handling another error (such as
// a logging failure during error reporting); the logging error may be associated
// with the original error.
//
// Examples of associated errors:
//   - In a file upload process, if both a network timeout and a disk write error occur
//     independently, each error can be associated with the other.
//   - If an error occurs while cleaning up resources after a primary error, the cleanup
//     error can be associated with the primary error.
//   - In a distributed system, if multiple nodes report errors for the same operation,
//     each node's error can be associated with the others.
//
// The returned slice may be empty or nil if there are no associated errors.
type ErrorAssociated interface {
	error

	// ErrorAssociated returns the associated errors for this error.
	// The returned slice may be nil or empty if there are no associated errors.
	ErrorAssociated() []error
}

// Associated returns the associated errors for the provided error, if any.
//
// This function attempts to extract the associated errors of the error as follows:
//  1. If the error implements ErrorAssociated, it returns the result of ErrorAssociated().
//  2. If the error is nil or does not implement ErrorAssociated, it returns nil.
//
// Associated errors are errors that are related to the current error, but are not
// direct causes. The returned slice may be nil or empty if there are no associated errors.
func Associated(err error) []error {
	if err == nil {
		return nil
	}

	if associated, ok := err.(ErrorAssociated); ok {
		return associated.ErrorAssociated()
	}

	return nil
}

// WithAssociated returns a new error with the given associated errors attached.
//
// This function takes an existing error and one or more associated errors, and returns a new error
// that includes the provided associated errors. If no associated errors are provided, it returns the original error.
// If the provided error is nil, it returns nil.
//
// Associated errors are errors that are related to the current error, but are not direct causes.
// This is useful for representing additional context, such as errors that occurred during cleanup,
// in parallel operations, or while handling the original error.
//
// Example:
//
//	err := fail.WithAssociated(primaryErr, cleanupErr, loggingErr)
//
// The returned error will implement the ErrorAssociated interface, and the associated errors
// can be retrieved using the fail.Associated(err) function.
//
// Parameters:
//   - err: The original error to which associated errors will be attached.
//   - associated: One or more errors to associate with the original error.
//
// Returns:
//   - A new error with the associated errors attached, or nil if err is nil. If no associated errors are provided, returns the original error.
func WithAssociated(err error, associated ...error) error {
	if err == nil {
		return nil
	}

	if len(associated) == 0 {
		return err
	}

	return From(err).Associate(associated...).asFail()
}
