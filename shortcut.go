package fail

import "context"

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
