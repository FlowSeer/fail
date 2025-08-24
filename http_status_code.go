package fail

// ErrorHttpStatusCode is an error type that provides an associated HTTP status code.
//
// Implementations of this interface should return a valid HTTP status code (such as 404, 500)
// to indicate the nature of the error in HTTP responses. If no specific status code is set,
// ErrorHttpStatusCode() should return 500 by default.
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
	// If no status code is explicitly set, implementations should return 500.
	ErrorHttpStatusCode() int
}

// HttpStatusCode returns the HTTP status code for the provided error.
//
// This function determines the HTTP status code as follows:
//  1. If err is nil, it returns 200 (success).
//  2. If err implements ErrorHttpStatusCode, it returns the result of ErrorHttpStatusCode().
//  3. Otherwise, it recursively examines the direct causes of err (using Causes(err)).
//     If any cause implements ErrorHttpStatusCode, it returns the maximum status code found among them.
//  4. If no status code is found, it returns 500 (default internal server error).
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

	maxHttpStatusCode := 500
	for _, cause := range Causes(err) {
		if httpStatusCode, ok := cause.(ErrorHttpStatusCode); ok {
			if httpStatusCode.ErrorHttpStatusCode() > maxHttpStatusCode {
				maxHttpStatusCode = httpStatusCode.ErrorHttpStatusCode()
			}
		}
	}

	return maxHttpStatusCode
}
