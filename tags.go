package fail

import (
	"maps"
	"slices"
)

// ErrorTags is an error type that provides a set of tags associated with the error.
//
// Implementations of this interface should return a slice of unique strings representing
// tags that describe or categorize the error. Tags can be used for filtering, logging,
// or error introspection. The returned slice must not contain duplicate tags, and may
// be empty or nil if there are no tags. The returned slice should be a copy, not a reference
// to internal state, to prevent callers from mutating the error's internal tags.
//
// Example usage:
//
//	type MyError struct{}
//	func (e *MyError) Error() string { return "something went wrong" }
//	func (e *MyError) ErrorTags() []string { return []string{"network", "timeout"} }
//
//	err := &MyError{}
//	tags := fail.Tags(err) // returns []string{"network", "timeout"}
type ErrorTags interface {
	error

	// ErrorTags returns the tags associated with this error.
	//
	// The returned slice must contain unique tag strings, and may be empty or nil if there are no tags.
	// The returned slice should be a copy, not a reference to internal state.
	ErrorTags() []string
}

// Tags returns the tags associated with the provided error, if any.
//
// This function attempts to extract tags from the error as follows:
//  1. If err is nil, it returns nil.
//  2. If err implements ErrorTags, it returns a deduplicated slice of tags from ErrorTags().
//     The returned slice is always a copy and safe for the caller to modify.
//  3. Otherwise, it returns nil.
//
// The returned slice may be nil or empty if there are no tags. The slice is always deduplicated
// and safe for the caller to modify.
func Tags(err error) []string {
	if err == nil {
		return nil
	}

	if tags, ok := err.(ErrorTags); ok {
		// We deduplicate the tags by using a map again to avoid duplicates from faulty implementations.
		tagsUniq := make(map[string]struct{})
		for _, t := range tags.ErrorTags() {
			tagsUniq[t] = struct{}{}
		}

		return slices.Collect(maps.Keys(tagsUniq))
	}

	return nil
}
