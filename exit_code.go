package fail

// DefaultExitCode is the default exit code to use when no specific exit code is set.
const DefaultExitCode = 1

// ErrorExitCode is an error type that provides a program exit code.
//
// Implementations of this interface should return a non-zero exit code to indicate failure.
// If no specific exit code is set, ErrorExitCode() should return DefaultExitCode by default.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorExitCode() int { return 42 }
//
//	err := &MyError{}
//	code := fail.ExitCode(err) // returns 42
type ErrorExitCode interface {
	error

	// ErrorExitCode returns the program exit code associated with this error.
	//
	// This value is intended to be used as the process exit status.
	// If no exit code is explicitly set, implementations should return DefaultExitCode.
	ErrorExitCode() int
}

// ExitCode returns the program exit code for the provided error.
//
// This function determines the exit code as follows:
//  1. If err is nil, it returns 0 (success).
//  2. If err implements ErrorExitCode, it returns the result of ErrorExitCode().
//  3. Otherwise, it recursively examines the direct causes of err (using Causes(err)).
//     If any cause implements ErrorExitCode, it returns the maximum exit code found among them.
//  4. If no exit code is found, it returns DefaultExitCode.
//
// This allows error types to specify custom exit codes, and for composed/multi-cause errors
// to propagate the most severe exit code.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	if exitCode, ok := err.(ErrorExitCode); ok {
		return exitCode.ErrorExitCode()
	}

	maxExitCode := DefaultExitCode
	for _, cause := range Causes(err) {
		if exitCode, ok := cause.(ErrorExitCode); ok {
			if exitCode.ErrorExitCode() > maxExitCode {
				maxExitCode = exitCode.ErrorExitCode()
			}
		}
	}

	return maxExitCode
}

// WithExitCode returns a new error with the specified program exit code attached.
//
// This function takes an existing error and an integer exit code, and returns a new error
// that includes the provided exit code. If the provided error is nil, it returns nil.
// If the exit code is less than or equal to zero, the original error is returned unchanged.
//
// The returned error will implement the ErrorExitCode interface, and the exit code can be
// retrieved using the fail.ExitCode function.
//
// Example:
//
//	err := fail.WithExitCode(primaryErr, 2)
//
// The returned error will have the exit code attached, and it can be accessed via
// fail.ExitCode(err).
//
// Parameters:
//   - err: The original error to which the exit code will be attached.
//   - exitCode: The integer exit code to associate with the error.
//
// Returns:
//   - A new error with the exit code attached, or nil if err is nil. If exitCode <= 0, returns the original error.
func WithExitCode(err error, exitCode int) error {
	if err == nil {
		return nil
	}

	if exitCode <= 0 {
		return err
	}

	return From(err).ExitCode(exitCode).asFail()
}
