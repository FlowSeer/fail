package fail

// DefaultHttpStatusCode is the default HTTP status code to use when no specific status code is set.
const DefaultHttpStatusCode = 500

// ErrorHttpStatusCode is an error type that provides an associated HTTP status code.
//
// Implementations of this interface should return a valid HTTP status code (such as 404, 500)
// to indicate the nature of the error in HTTP responses. If no specific status code is set,
// ErrorHttpStatusCode() should return DefaultHttpStatusCode by default.
//
// Example usage:
//
//	type MyHTTPError struct{}
//	func (e *MyHTTPError) Error() string { return "not found" }
//	func (e *MyHTTPError) ErrorHttpStatusCode() int { return 404 }
//
//	err := &MyHTTPError{}
//	code := fail.HttpStatusCode(err) // returns 404
type ErrorHttpStatusCode interface {
	error

	// ErrorHttpStatusCode returns the HTTP status code associated with this error.
	//
	// This value is intended to be used as the HTTP response status.
	// If no status code is explicitly set, implementations should return DefaultHttpStatusCode.
	ErrorHttpStatusCode() int
}

// HttpStatusCode returns the HTTP status code for the provided error.
//
// This function determines the HTTP status code as follows:
//  1. If err is nil, it returns 200 (success).
//  2. If err implements ErrorHttpStatusCode, it returns the result of ErrorHttpStatusCode().
//  3. Otherwise, it recursively examines the direct causes of err (using Causes(err)).
//     If any cause implements ErrorHttpStatusCode, it returns the maximum status code found among them.
//  4. If no status code is found, it returns DefaultHttpStatusCode.
//
// This allows error types to specify custom HTTP status codes, and for composed/multi-cause errors
// to propagate the most severe status code.
func HttpStatusCode(err error) int {
	if err == nil {
		return 200
	}

	if httpStatusCode, ok := err.(ErrorHttpStatusCode); ok {
		return httpStatusCode.ErrorHttpStatusCode()
	}

	maxHttpStatusCode := DefaultHttpStatusCode
	for _, cause := range Causes(err) {
		if httpStatusCode, ok := cause.(ErrorHttpStatusCode); ok {
			if httpStatusCode.ErrorHttpStatusCode() > maxHttpStatusCode {
				maxHttpStatusCode = httpStatusCode.ErrorHttpStatusCode()
			}
		}
	}

	return maxHttpStatusCode
}

// WithHttpStatusCode returns a new error with the specified HTTP status code attached.
//
// This function takes an existing error and an integer HTTP status code, and returns a new error
// that includes the provided status code. If the provided error is nil, it returns nil.
// If the HTTP status code is less than 400, the original error is returned unchanged.
//
// The returned error will implement the ErrorHttpStatusCode interface, and the status code can be
// retrieved using the fail.HttpStatusCode function.
//
// Example:
//
//	err := fail.WithHttpStatusCode(primaryErr, 404)
//
// The returned error will have the HTTP status code attached, and it can be accessed via
// fail.HttpStatusCode(err).
//
// Parameters:
//   - err: The original error to which the HTTP status code will be attached.
//   - httpStatusCode: The integer HTTP status code to associate with the error.
//
// Returns:
//   - A new error with the HTTP status code attached, or nil if err is nil. If httpStatusCode < 400, returns the original error.
func WithHttpStatusCode(err error, httpStatusCode int) error {
	if err == nil {
		return nil
	}

	if httpStatusCode < 400 {
		return err
	}

	return From(err).HttpStatusCode(httpStatusCode).asFail()
}
