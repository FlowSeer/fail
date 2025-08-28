package fail

import "context"

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

// WithAttributes returns a new error with the given attributes attached.
//
// This function takes an existing error and a map of attributes, and returns a new error
// that includes the provided attributes. If the provided error is nil, it returns nil.
// If the attrs map is empty or nil, the original error is returned unchanged.
//
// The returned error will implement the ErrorAttributes interface, and the attributes
// can be retrieved using the fail.Attributes function.
//
// Example:
//
//	err := fail.WithAttributes(primaryErr, map[string]any{"foo": 42, "bar": "baz"})
//
// The returned error will have the attributes attached, and they can be accessed via
// fail.Attributes(err).
//
// Parameters:
//   - err: The original error to which attributes will be attached.
//   - attrs: A map of key-value attributes to associate with the error.
//
// Returns:
//   - A new error with the attributes attached, or nil if err is nil. If attrs is empty or nil, returns the original error.
func WithAttributes(err error, attrs map[string]any) error {
	if err == nil {
		return nil
	}

	if len(attrs) == 0 {
		return err
	}

	return From(err).AttributeMap(attrs).asFail()
}

// attributesContextKey is an unexported type used as the key for storing
// and retrieving error attributes in a context.Context.
type attributesContextKey struct{}

// ContextWithAttributes returns a new context.Context that carries the provided
// error attributes map. If attributes are already set in the context, they are overwritten
// with the new value.
//
// Example usage:
//
//	ctx := ContextWithAttributes(context.Background(), map[string]any{"foo": 42})
func ContextWithAttributes(ctx context.Context, attrs map[string]any) context.Context {
	return context.WithValue(ctx, attributesContextKey{}, attrs)
}

// ContextAddAttributes returns a new context.Context with the provided attributes merged
// into any existing attributes in the context. If no attributes are present, it behaves like ContextWithAttributes.
// If the same key exists in both the existing and new attributes, the value from attrs overwrites the existing value.
//
// Example usage:
//
//	ctx := ContextAddAttributes(ctx, map[string]any{"bar": "baz"})
func ContextAddAttributes(ctx context.Context, attrs map[string]any) context.Context {
	existingAttrs := AttributesFromContext(ctx)
	merged := make(map[string]any, len(existingAttrs)+len(attrs))
	for k, v := range existingAttrs {
		merged[k] = v
	}
	for k, v := range attrs {
		merged[k] = v
	}
	return ContextWithAttributes(ctx, merged)
}

// AttributesFromContext extracts the error attributes map from the provided context.
// If no attributes are set in the context, AttributesFromContext returns nil.
//
// Example usage:
//
//	attrs := AttributesFromContext(ctx)
func AttributesFromContext(ctx context.Context) map[string]any {
	attrs, ok := ctx.Value(attributesContextKey{}).(map[string]any)
	if !ok {
		return nil
	}
	return attrs
}
