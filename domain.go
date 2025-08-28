package fail

import "context"

// Predefined error domain constants for categorizing errors by their source or type.
//
// These constants can be used to tag errors with a specific domain, enabling
// easier grouping, filtering, and handling of errors throughout your application.
//
// Common domains include network issues, configuration problems, database errors,
// validation failures, authentication issues, rate limiting, and more. You can
// extend this list as needed to fit your application's requirements.
const (
	// DomainUnknown represents an unknown or uncategorized error domain.
	DomainUnknown = "unknown"
	// DomainNetwork represents errors related to network connectivity or communication.
	DomainNetwork = "network"
	// DomainConfig represents errors related to configuration, such as missing or invalid settings.
	DomainConfig = "config"
	// DomainDatabase represents errors originating from database operations.
	DomainDatabase = "database"
	// DomainValidation represents errors due to validation failures (e.g., invalid input).
	DomainValidation = "validation"
	// DomainAuth represents authentication or authorization errors.
	DomainAuth = "auth"
	// DomainRateLimit represents errors caused by exceeding rate limits.
	DomainRateLimit = "ratelimit"
	// DomainIO represents errors related to input/output operations (e.g., file system).
	DomainIO = "io"
	// DomainTimeout represents errors caused by operation timeouts.
	DomainTimeout = "timeout"
	// DomainDependency represents errors from external dependencies or services.
	DomainDependency = "dependency"
	// DomainInternal represents internal application errors not exposed to users.
	DomainInternal = "internal"
	// DomainAPI represents errors related to API usage or responses.
	DomainAPI = "api"
)

// ErrorDomain is an interface for categorizing errors by domain.
//
// Implementations of ErrorDomain should return a unique, short string
// that identifies the domain or category to which the error belongs.
// This enables grouping, filtering, or handling errors based on their domain.
//
// Example domains might include "network", "database", or "validation".
type ErrorDomain interface {
	error

	// ErrorDomain returns the unique domain name of the error as a string.
	// This value should be consistent and unique for each error domain.
	ErrorDomain() string
}

// Domain returns the domain name of the given error if it implements the ErrorDomain interface.
//
// If the error is nil, Domain returns an empty string. If the error implements
// ErrorDomain, Domain returns the result of its ErrorDomain() method. Otherwise,
// it returns DomainUnknown.
//
// Example usage:
//
//	type MyError struct{}
//
//	func (MyError) Error() string { return "something went wrong" }
//	func (MyError) ErrorDomain() string { return "custom" }
//
//	err := MyError{}
//	domain := fail.Domain(err) // domain == "custom"
func Domain(err error) string {
	if err == nil {
		return ""
	}

	if domain, ok := err.(ErrorDomain); ok {
		return domain.ErrorDomain()
	}

	return DomainUnknown
}

// WithDomain returns a new error with the specified domain attached.
//
// This function takes an existing error and a domain string, and returns a new error
// that includes the provided domain. If the provided error is nil, it returns nil.
// If the domain string is empty, the original error is returned unchanged.
//
// The returned error will implement the ErrorDomain interface, and the domain can be
// retrieved using the fail.Domain function.
//
// Example:
//
//	err := fail.WithDomain(primaryErr, fail.DomainDatabase)
//
// The returned error will have the domain attached, and it can be accessed via
// fail.Domain(err).
//
// Parameters:
//   - err: The original error to which the domain will be attached.
//   - domain: The domain string to associate with the error.
//
// Returns:
//   - A new error with the domain attached, or nil if err is nil. If domain is empty, returns the original error.
func WithDomain(err error, domain string) error {
	if err == nil {
		return nil
	}

	if domain == "" {
		return err
	}

	return From(err).Domain(domain).asFail()
}

// domainContextKey is an unexported type used as the key for storing
// and retrieving the error domain value in a context.Context.
type domainContextKey struct{}

// ContextWithDomain returns a new context.Context that carries the provided
// error domain string. If a domain is already set in the context, it is overwritten
// with the new value. This allows the error domain to be propagated
// through request or operation lifecycles via context.
//
// Example usage:
//
//	ctx := ContextWithDomain(context.Background(), DomainValidation)
func ContextWithDomain(ctx context.Context, domain string) context.Context {
	// context.WithValue always overwrites the value for the key if it already exists.
	return context.WithValue(ctx, domainContextKey{}, domain)
}

// DomainFromContext extracts the error domain string from the provided context.
// If no domain is set in the context, DomainUnknown is returned.
//
// Example usage:
//
//	domain := DomainFromContext(ctx)
func DomainFromContext(ctx context.Context) string {
	domain, ok := ctx.Value(domainContextKey{}).(string)
	if !ok {
		return DomainUnknown
	}
	return domain
}
