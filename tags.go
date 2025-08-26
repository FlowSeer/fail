package fail

import (
	"context"
	"maps"
	"slices"
)

// Predefined error domain constants for categorizing errors by their source or type.
//
// These constants can be used to tag errors with a specific domain, enabling
// easier grouping, filtering, and handling of errors throughout your application.
//
// Common domains include network issues, configuration problems, database errors,
// validation failures, authentication issues, rate limiting, and more. You can
// extend this list as needed to fit your application's requirements.
const (
	// TagNetwork represents errors related to network connectivity or communication.
	TagNetwork = DomainNetwork
	// TagConfig represents errors related to configuration, such as missing or invalid settings.
	TagConfig = DomainConfig
	// TagDatabase represents errors originating from database operations.
	TagDatabase = DomainDatabase
	// TagValidation represents errors due to validation failures (e.g., invalid input).
	TagValidation = DomainValidation
	// TagAuth represents authentication or authorization errors.
	TagAuth = DomainAuth
	// TagRateLimit represents errors caused by exceeding rate limits.
	TagRateLimit = DomainRateLimit
	// TagIO represents errors related to input/output operations (e.g., file system).
	TagIO = DomainIO
	// TagTimeout represents errors caused by operation timeouts.
	TagTimeout = DomainTimeout
	// TagDependency represents errors from external dependencies or services.
	TagDependency = DomainDependency
	// TagInternal represents internal application errors not exposed to users.
	TagInternal = DomainInternal
	// TagAPI represents errors related to API usage or responses.
	TagAPI = DomainAPI
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

// tagsContextKey is an unexported type used as the key for storing
// and retrieving error tags in a context.Context.
type tagsContextKey struct{}

// ContextWithTags returns a new context.Context that carries the provided
// error tags slice. If tags are already set in the context, they are overwritten
// with the new value. Use ContextAddTags to add tags without overwriting existing tags.
//
// Example usage:
//
//	ctx := ContextWithTags(context.Background(), []string{"network", "timeout"})
func ContextWithTags(ctx context.Context, tags []string) context.Context {
	return context.WithValue(ctx, tagsContextKey{}, tags)
}

// ContextAddTags returns a new context.Context with the provided tags appended
// to any existing tags in the context. If no tags are present, it behaves like ContextWithTags.
// The resulting tags slice may contain duplicates.
//
// Example usage:
//
//	ctx := ContextAddTags(ctx, []string{"database"})
func ContextAddTags(ctx context.Context, tags []string) context.Context {
	existingTags := TagsFromContext(ctx)
	return ContextWithTags(ctx, append(existingTags, tags...))
}

// TagsFromContext extracts the error tags slice from the provided context.
// If no tags are set in the context, TagsFromContext returns nil.
//
// Example usage:
//
//	tags := TagsFromContext(ctx)
func TagsFromContext(ctx context.Context) []string {
	tags, ok := ctx.Value(tagsContextKey{}).([]string)
	if !ok {
		return nil
	}
	return tags
}
