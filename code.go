package fail

// DefaultErrorCode is the fallback error code returned by Code if no specific code is set
// and no code is found in any error in the chain.
const DefaultErrorCode = "ERR_UNKNOWN"

// ErrorCode is an error type that provides a canonical, programmatic error code.
//
// Implementations of this interface should return a stable, concise string code
// that uniquely identifies the type or category of the error. The code must not contain
// whitespace or special characters—only letters, numbers, and underscores are allowed.
// This code is intended for programmatic use (e.g., logs, diagnostics, error handling),
// and must not include details from wrapped errors or underlying causes.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorCode() string { return "ERR1407" }
//
//	err := &MyError{}
//	code := fail.Code(err) // returns "ERR1407"
type ErrorCode interface {
	error

	// ErrorCode returns the canonical error code associated with this error.
	// The code must be a simple string identifying the error type, using only
	// letters, numbers, and underscores. No whitespace or special characters.
	ErrorCode() string
}

// Code returns the canonical programmatic error code for the provided error.
//
// This function determines the error code as follows:
//  1. If err is nil, it returns the empty string.
//  2. If err implements ErrorCode, it returns the result of ErrorCode().
//  3. Otherwise, it recursively examines the direct causes of err (using Causes(err)).
//     If any cause implements ErrorCode, it returns the code from the cause with the highest ExitCode.
//     Otherwise, it returns the code from the cause with the first non-default code.
//  4. If no code is found, it returns DefaultErrorCode.
//
// This allows error types to specify custom error codes, and for composed/multi-cause errors
// to propagate the code from the most severe cause (as determined by ExitCode).
func Code(err error) string {
	if err == nil {
		return ""
	}

	// If the error itself implements ErrorCode, return its code.
	if code, ok := err.(ErrorCode); ok {
		return code.ErrorCode()
	}

	// Otherwise, check causes and return the code from the cause with the highest exit code.
	maxCode := DefaultErrorCode
	maxExitCode := 0
	for _, cause := range Causes(err) {
		causeExitCode := ExitCode(cause)
		causeCode := Code(cause)

		// Prefer the code from the cause with the highest exit code.
		if causeExitCode > maxExitCode {
			maxExitCode = causeExitCode
			maxCode = causeCode
		} else if maxCode == DefaultErrorCode && causeCode != DefaultErrorCode {
			// If no better code has been found yet, use the first non-default code.
			maxCode = causeCode
		}
	}

	return maxCode
}
