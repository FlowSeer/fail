package fail

// Fatal prints the provided error to standard output and exits the program with a non-zero exit code.
// If the error is nil, it does nothing.
//
// Example:
//
//	fail.Fatal(err)
func Fatal(err error) {
	if err == nil {
		return
	}

	PrintPretty(err)
	Exit(err)
}

// FatalMsg prints the provided message to standard output and exits the program with a non-zero exit code.
// If the message is empty, it does nothing.
//
// Example:
//
//	fail.FatalMsg("failed to perform operation")
func FatalMsg(msg string) {
	if msg == "" {
		return
	}

	Fatal(Msg(msg))
}

// Fatalf prints the provided formatted message to standard output and exits the program with a non-zero exit code.
// If the message is empty, it does nothing.
//
// Example:
//
//	fail.Fatalf("failed to perform operation: %s", err)
func Fatalf(format string, args ...interface{}) {
	Fatal(Msgf(format, args...))
}
