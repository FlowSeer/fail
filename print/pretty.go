package print

import "fail"

// Pretty returns a human-readable string representation of the provided error.
//
// This function uses the default PrettyPrinter to format the error. It is suitable
// for displaying errors in logs, user interfaces, or diagnostics where a readable
// format is desired.
//
// Example:
//
//	err := fail.New().Msg("something went wrong")
//	out := print.Pretty(err)
func Pretty(err error) string {
	return PrettyPrinter().Print(err)
}

// PrettyPrinter returns a Printer that formats errors in a human-readable way.
//
// The returned Printer uses the provided PrinterOptions to control which fields
// are included in the output, such as causes, associated errors, codes, tags, etc.
// This is useful for customizing error output for logs or user interfaces.
//
// Example:
//
//	printer := print.PrettyPrinter(print.WithoutColor())
//	out := printer.Print(err)
func PrettyPrinter(opts ...PrinterOption) Printer {
	return PrinterFunc(func(err error) string {
		return printPretty(err, opts...)
	})
}

// printPretty formats the provided error as a human-readable string according to the given PrinterOptions.
//
// This is an internal helper used by PrettyPrinter and Pretty. Currently, it returns only the error message.
// In the future, it may be extended to include more error metadata.
func printPretty(err error, opts ...PrinterOption) string {
	return fail.Message(err)
}
