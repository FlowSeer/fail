package fail

// ErrorExitCode is an error type that provides a program exit code.
//
// Implementations of this interface should return a non-zero exit code to indicate failure.
// If no specific exit code is set, ErrorExitCode() should return 1 by default.
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
	// If no exit code is explicitly set, implementations should return 1.
	ErrorExitCode() int
}

// ExitCode returns the program exit code for the provided error.
//
// This function determines the exit code as follows:
//  1. If err is nil, it returns 0 (success).
//  2. If err implements ErrorExitCode, it returns the result of ErrorExitCode().
//  3. Otherwise, it recursively examines the direct causes of err (using Causes(err)).
//     If any cause implements ErrorExitCode, it returns the maximum exit code found among them.
//  4. If no exit code is found, it returns 1 (default failure exit code).
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

	maxExitCode := 1
	for _, cause := range Causes(err) {
		if exitCode, ok := cause.(ErrorExitCode); ok {
			if exitCode.ErrorExitCode() > maxExitCode {
				maxExitCode = exitCode.ErrorExitCode()
			}
		}
	}

	return maxExitCode
}
