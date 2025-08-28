package fail

// ErrorCauses is an error type that provides a list of underlying causes for an error.
//
// Implementations of this interface should return a slice of errors representing the
// direct causes of the current error. This is useful for error introspection, error
// chaining, and for building error trees or graphs. The returned slice may be empty
// or nil if there are no underlying causes.
//
// Example usage:
//
//	type MyError struct {
//	    causes []error
//	}
//
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorCauses() []error { return e.causes }
//
//	err := &MyError{causes: []error{io.EOF, io.ErrClosedPipe}}
//	causes := fail.Causes(err) // returns []error{io.EOF, io.ErrClosedPipe}
type ErrorCauses interface {
	error

	// ErrorCauses returns the direct underlying causes of this error.
	// The returned slice may be nil or empty if there are no causes.
	ErrorCauses() []error
}

// Causes returns the direct underlying causes of the provided error, if any.
//
// This function attempts to extract the causes of the error in the following order:
//  1. If the error implements ErrorCauses, it returns the result of ErrorCauses().
//  2. If the error implements Unwrap() []error, it returns the result of Unwrap().
//  3. If the error implements Unwrap() error, it returns a single-element slice containing the result of Unwrap().
//  4. If the error implements Cause() error (as in github.com/pkg/errors), it returns a single-element slice containing the result of Cause().
//  5. If none of the above, or if err is nil, it returns nil.
//
// The returned slice may be nil or empty if there are no causes.
func Causes(err error) []error {
	if err == nil {
		return nil
	}

	// Check if the error implements ErrorCauses.
	if causes, ok := err.(ErrorCauses); ok {
		return causes.ErrorCauses()
	}

	// Check if the error implements Unwrap() []error (Go 1.20+ multi-error).
	if unwrapSlice, ok := err.(interface{ Unwrap() []error }); ok {
		return unwrapSlice.Unwrap()
	}

	// Check if the error implements Unwrap() error (Go 1.13+).
	if unwrap, ok := err.(interface{ Unwrap() error }); ok {
		return []error{unwrap.Unwrap()}
	}

	// Check for the legacy Cause() error method (e.g., github.com/pkg/errors).
	if cause, ok := err.(interface{ Cause() error }); ok {
		return []error{cause.Cause()}
	}

	return nil
}

// WithCauses returns a new error with the specified direct causes attached.
//
// This function takes an existing error and one or more direct causes, and returns a new error
// that includes the provided causes. If the provided error is nil, it returns nil. If no causes
// are provided, the original error is returned unchanged.
//
// Causes are errors that directly led to the current error. This is useful for representing
// error chains or aggregating multiple underlying errors that contributed to the failure.
//
// Example:
//
//	err := fail.WithCauses(primaryErr, io.EOF, io.ErrClosedPipe)
//
// The returned error will implement the ErrorCauses interface, and the causes can be retrieved
// using fail.Causes(err).
//
// Parameters:
//   - err: The original error to which causes will be attached.
//   - causes: One or more errors to set as the direct causes of the error.
//
// Returns:
//   - A new error with the causes attached, or nil if err is nil. If no causes are provided, returns the original error.
func WithCauses(err error, causes ...error) error {
	if err == nil {
		return nil
	}

	if len(causes) == 0 {
		return err
	}

	return From(err).Cause(causes...).asFail()
}
