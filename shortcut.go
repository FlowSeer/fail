package fail

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

// Wrap creates a new Fail error with the given message and wraps the provided error as its cause.
//
// If err is nil, Wrap returns nil.
// This is a shortcut for fail.New().Cause(err).Msg(msg).
//
// Example:
//
//	err := fail.Wrap("failed to read file", io.EOF)
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return New().Cause(err).Msg(msg)
}

// Wrapf creates a new Fail error with a formatted message and wraps the provided error as its cause.
//
// This is a shortcut for fail.New().Cause(err).Msgf(format, args...).
//
// Example:
//
//	err := fail.Wrapf(io.EOF, "failed to read file %q", filename)
func Wrapf(err error, format string, args ...any) error {
	return New().Cause(err).Msgf(format, args...)
}

// WrapMany creates a new Fail error with the given message and wraps multiple errors as its causes.
//
// If errs is empty, WrapMany returns nil.
// This is a shortcut for fail.New().CauseSlice(errs).Msg(msg).
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
