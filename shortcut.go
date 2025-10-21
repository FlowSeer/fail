package fail

import (
	"context"
)

// Msg creates a new Fail error with the given developer-facing message.
//
// This is a shortcut for fail.New().Msg(msg).
// The returned error implements all fail.* error interfaces.
//
// Example:
//
//	err := fail.Msg("database connection failed")
func Msg(msg string) error {
	return New().Msg(msg)
}

// MsgC creates a new Fail error with the given developer-facing message and context.
//
// This is a shortcut for fail.NewC(ctx).Msg(msg).
// The returned error implements all fail.* error interfaces.
//
// Example:
//
//	err := fail.MsgC(ctx, "database connection failed")
func MsgC(ctx context.Context, msg string) error {
	return NewC(ctx).Msg(msg)
}

// Msgf creates a new Fail error with a formatted developer-facing message.
//
// This is a shortcut for fail.New().Msgf(format, args...).
// The returned error implements all fail.* error interfaces.
//
// Example:
//
//	err := fail.Msgf("failed to connect to database %s on port %d", "localhost", 5432)
func Msgf(format string, args ...any) error {
	return New().Msgf(format, args...)
}

// MsgCf creates a new Fail error with a formatted developer-facing message and context.
//
// This is a shortcut for fail.NewC(ctx).Msgf(format, args...).
// The returned error implements all fail.* error interfaces.
//
// Example:
//
//	err := fail.MsgCf(ctx, "failed to connect to database %s on port %d", "localhost", 5432)
func MsgCf(ctx context.Context, format string, args ...any) error {
	return NewC(ctx).Msgf(format, args...)
}

// Wrap returns a new Fail error with the given message, wrapping the provided error as its cause.
//
// If err is nil, Wrap returns nil.
// Equivalent to: fail.New().Cause(err).Msg(msg).
//
// Example:
//
//	err := fail.Wrap(io.EOF, "failed to read file")
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return New().Cause(err).Msg(msg)
}

// WrapResult calls the provided function fn, and if it returns a non-nil error,
// wraps that error using fail.Wrap with the given message. The result value is
// returned as-is. This helper is commonly used to propagate errors with context
// when working in functional or result-oriented code.
//
// Usage:
//
//	result, err := fail.WrapResult(doSomething, "failed operation")
//	if err != nil {
//	    handle(err) // err is wrapped with message "failed operation"
//	} else {
//	    use(result)
//	}
//
// It is especially helpful in higher-order patterns:
//
//	handle(fail.WrapResult(action, "failed some action"))
func WrapResult[T any](fn func() (T, error), msg string) (T, error) {
	res, err := fn()
	return res, Wrap(err, msg)
}

// WrapC creates a new Fail error with the given message, wrapping the provided error as its cause and context.
//
// If err is nil, WrapC returns nil.
// Equivalent to: fail.NewC(ctx).Cause(err).Msg(msg).
//
// Example:
//
//	err := fail.WrapC(ctx, io.EOF, "failed to read file")
func WrapC(ctx context.Context, err error, msg string) error {
	if err == nil {
		return nil
	}

	return NewC(ctx).Cause(err).Msg(msg)
}

// WrapCResult executes the provided function fn, and if it returns a non-nil error,
// wraps that error with additional context using the given message and context.
// The result value is returned as-is. This is useful when you need to propagate
// error context with a message and a context.Context.
//
// Usage:
//
//	val, err := fail.WrapCResult(ctx, someFunc, "failed to perform operation")
//	if err != nil {
//	    handle(err) // err is wrapped with message and context
//	} else {
//	    use(val)
//	}
//
// Equivalent to calling: result, err := fn(); err = fail.WrapC(ctx, err, msg)
func WrapCResult[T any](ctx context.Context, fn func() (T, error), msg string) (T, error) {
	res, err := fn()
	return res, WrapC(ctx, err, msg)
}

// Wrapf returns a new Fail error with a formatted message, wrapping the provided error as its cause.
//
// Equivalent to: fail.New().Cause(err).Msgf(format, args...).
//
// Example:
//
//	err := fail.Wrapf(io.EOF, "failed to read file %q", filename)
func Wrapf(err error, format string, args ...any) error {
	return New().Cause(err).Msgf(format, args...)
}

func WrapfResult[T any](fn func() (T, error), format string, args ...any) (T, error) {
	res, err := fn()
	return res, Wrapf(err, format, args...)
}

// WrapCf creates a new Fail error with a formatted message, wrapping the provided error as its cause and context.
//
// If err is nil, WrapCf returns nil.
// Equivalent to: fail.NewC(ctx).Cause(err).Msgf(format, args...).
//
// Example:
//
//	err := fail.WrapCf(ctx, io.EOF, "failed to read file %q", filename)
func WrapCf(ctx context.Context, err error, format string, args ...any) error {
	return NewC(ctx).Cause(err).Msgf(format, args...)
}

// WrapCfResult executes the provided function fn, and if it returns a non-nil error,
// wraps that error with a formatted message using the specified context and format string.
//
// This is useful when you want to add both context.Context and a formatted message to an error
// returned by a function, and still return the result value transparently.
//
// Usage:
//
//	val, err := fail.WrapCfResult(ctx, someFunc, "failed to process item %d", id)
//	if err != nil {
//	    handle(err) // err is wrapped with context and formatted message
//	} else {
//	    use(val)
//	}
//
// Equivalent to calling: result, err := fn(); err = fail.WrapCf(ctx, err, format, args...)
func WrapCfResult[T any](ctx context.Context, fn func() (T, error), format string, args ...any) (T, error) {
	res, err := fn()
	return res, WrapCf(ctx, err, format, args...)
}

// WrapMany returns a new Fail error with the given message, wrapping multiple errors as its causes.
//
// If errs is empty, WrapMany returns nil.
// Equivalent to: fail.New().CauseSlice(errs).Msg(msg).
//
// Example:
//
//	err := fail.WrapMany("multiple errors occurred", err1, err2, err3)
func WrapMany(msg string, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return New().CauseSlice(errs).Msg(msg)
}

// WrapManyC creates a new Fail error with the given message, wrapping multiple errors as its causes and context.
//
// If errs is empty, WrapManyC returns nil.
// Equivalent to: fail.NewC(ctx).CauseSlice(errs).Msg(msg).
//
// Example:
//
//	err := fail.WrapManyC(ctx, "multiple errors occurred", err1, err2, err3)
func WrapManyC(ctx context.Context, msg string, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return NewC(ctx).CauseSlice(errs).Msg(msg)
}

// WithContext adds information from the provided context to the error.
//
// If err is nil, WithContext returns nil.
func WithContext(err error, ctx context.Context) error {
	if err == nil {
		return nil
	}

	return From(err).Context(ctx).asFail()
}
