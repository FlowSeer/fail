package fail

// ErrorAttributes is an error type that provides a set of key-value attributes associated with the error.
//
// Implementations of this interface should return a map of attributes that describe or provide
// additional context for the error. The returned map must be a clone (not the internal map itself)
// to prevent callers from mutating the error's internal state. The map must not be nil; if there
// are no attributes, an empty map must be returned.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorAttributes() map[string]any { return map[string]any{"foo": 42, "bar": "baz"} }
//
//	err := &MyError{}
//	attrs := fail.Attributes(err) // returns map[string]any{"foo": 42, "bar": "baz"}
type ErrorAttributes interface {
	error

	// ErrorAttributes returns the attributes associated with this error.
	//
	// The returned map must be a clone of the internal attributes map, so that modifications
	// to the returned map do not affect the error's internal state. The map must not be nil;
	// if there are no attributes, an empty map must be returned.
	ErrorAttributes() map[string]any
}

// Attributes returns the attributes associated with the provided error, if any.
//
// This function attempts to extract attributes from the error as follows:
//  1. If err is nil, it returns an empty map.
//  2. If err implements ErrorAttributes which returns a non-nil map, it returns the result of ErrorAttributes().
//  3. Otherwise, it returns an empty map.
//
// The returned map is always non-nil and safe for the caller to modify. If there are no attributes, an empty map is returned.
func Attributes(err error) map[string]any {
	if err == nil {
		return map[string]any{}
	}
	if attrs, ok := err.(ErrorAttributes); ok {
		if attrsMap := attrs.ErrorAttributes(); attrsMap != nil {
			return attrsMap
		}
	}
	return map[string]any{}
}
