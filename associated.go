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
