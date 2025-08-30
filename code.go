package fail

// Error code constants for canonical programmatic error codes.
const (
	// ErrCodeUnspecified is the default error code for unknown or unspecified errors.
	ErrCodeUnspecified = "ERR_UNKNOWN"
	// ErrCodeUnknown is an alias for ErrCodeUnspecified.
	ErrCodeUnknown = ErrCodeUnspecified

	// Validation errors
	// ErrCodeValidation indicates a general validation failure.
	ErrCodeValidation = "ERR_VALIDATION"
	// ErrCodeInvalidInput indicates that input data is invalid.
	ErrCodeInvalidInput = "ERR_INVALID_INPUT"
	// ErrCodeMissingRequired indicates a required value is missing.
	ErrCodeMissingRequired = "ERR_MISSING_REQUIRED"
	// ErrCodeInvalidFormat indicates data is in an invalid format.
	ErrCodeInvalidFormat = "ERR_INVALID_FORMAT"
	// ErrCodeOutOfRange indicates a value is outside the allowed range.
	ErrCodeOutOfRange = "ERR_OUT_OF_RANGE"

	// Authentication and authorization errors
	// ErrCodeUnauthorized indicates the user is not authorized.
	ErrCodeUnauthorized = "ERR_UNAUTHORIZED"
	// ErrCodeForbidden indicates access is forbidden.
	ErrCodeForbidden = "ERR_FORBIDDEN"
	// ErrCodeAuthentication indicates a general authentication failure.
	ErrCodeAuthentication = "ERR_AUTHENTICATION"
	// ErrCodeTokenExpired indicates an authentication token has expired.
	ErrCodeTokenExpired = "ERR_TOKEN_EXPIRED"
	// ErrCodeInvalidToken indicates an authentication token is invalid.
	ErrCodeInvalidToken = "ERR_INVALID_TOKEN"

	// Resource errors
	// ErrCodeNotFound indicates a requested resource was not found.
	ErrCodeNotFound = "ERR_NOT_FOUND"
	// ErrCodeAlreadyExists indicates a resource already exists.
	ErrCodeAlreadyExists = "ERR_ALREADY_EXISTS"
	// ErrCodeConflict indicates a resource conflict.
	ErrCodeConflict = "ERR_CONFLICT"
	// ErrCodeResourceGone indicates a resource is no longer available.
	ErrCodeResourceGone = "ERR_RESOURCE_GONE"

	// Network and communication errors
	// ErrCodeNetwork indicates a general network error.
	ErrCodeNetwork = "ERR_NETWORK"
	// ErrCodeTimeout indicates a timeout occurred.
	ErrCodeTimeout = "ERR_TIMEOUT"
	// ErrCodeConnection indicates a connection error.
	ErrCodeConnection = "ERR_CONNECTION"
	// ErrCodeUnreachable indicates a resource or service is unreachable.
	ErrCodeUnreachable = "ERR_UNREACHABLE"

	// System and infrastructure errors
	// ErrCodeInternal indicates an internal system error.
	ErrCodeInternal = "ERR_INTERNAL"
	// ErrCodeServiceUnavailable indicates a service is unavailable.
	ErrCodeServiceUnavailable = "ERR_SERVICE_UNAVAILABLE"
	// ErrCodeDatabase indicates a database error.
	ErrCodeDatabase = "ERR_DATABASE"
	// ErrCodeStorage indicates a storage error.
	ErrCodeStorage = "ERR_STORAGE"
	// ErrCodeConfiguration indicates a configuration error.
	ErrCodeConfiguration = "ERR_CONFIGURATION"

	// Business logic errors
	// ErrCodeBusinessRule indicates a business rule violation.
	ErrCodeBusinessRule = "ERR_BUSINESS_RULE"
	// ErrCodeQuotaExceeded indicates a quota has been exceeded.
	ErrCodeQuotaExceeded = "ERR_QUOTA_EXCEEDED"
	// ErrCodeRateLimited indicates a rate limit has been exceeded.
	ErrCodeRateLimited = "ERR_RATE_LIMITED"
	// ErrCodeMaintenance indicates the system is in maintenance mode.
	ErrCodeMaintenance = "ERR_MAINTENANCE"
)

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
	maxCode := ErrCodeUnspecified
	maxExitCode := 0
	for _, cause := range Causes(err) {
		causeExitCode := ExitCode(cause)
		causeCode := Code(cause)

		// Prefer the code from the cause with the highest exit code.
		if causeExitCode > maxExitCode {
			maxExitCode = causeExitCode
			maxCode = causeCode
		} else if maxCode == ErrCodeUnspecified && causeCode != ErrCodeUnspecified {
			// If no better code has been found yet, use the first non-default code.
			maxCode = causeCode
		}
	}

	return maxCode
}

func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}

	return From(err).Code(code).asFail()
}
