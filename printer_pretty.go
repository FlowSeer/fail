package fail

import (
	"strings"
)

// PrintPretty prints a human-readable string representation of the provided error to standard output.
//
// This function uses the default PrettyPrinter to format the error. It is suitable
// for displaying errors in logs, user interfaces, or diagnostics where a readable
// format is desired.
//
// Example:
//
//	err := fail.New().Msg("something went wrong")
//	print.PrintPretty(err)
func PrintPretty(err error, opts ...PrinterOption) {
	println(PrintsPretty(err, opts...))
}

// PrintsPretty returns a human-readable string representation of the provided error.
//
// This function uses the default PrettyPrinter to format the error. It is suitable
// for displaying errors in logs, user interfaces, or diagnostics where a readable
// format is desired.
//
// Example:
//
//	err := fail.New().Msg("something went wrong")
//	out := print.PrintsPretty(err)
func PrintsPretty(err error, opts ...PrinterOption) string {
	return PrettyPrinter(opts...).Print(err)
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
	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	return PrinterFunc(func(err error) string {
		sb := strings.Builder{}
		printPretty(&sb, 0, err, o)

		return sb.String()
	})
}

// printPretty formats the provided error as a human-readable string according to the given PrinterOptions.
//
// This is an internal helper used by PrettyPrinter and PrintPretty. Currently, it returns only the error message.
// In the future, it may be extended to include more error metadata.
// TODO: improve logging
func printPretty(sb *strings.Builder, depth int, err error, opts PrinterOptions) {
	sb.WriteString(strings.Repeat("  ", depth) + Message(err))

	if opts.Causes && (opts.CauseDepth == 0 || depth <= opts.CauseDepth) {
		for _, cause := range Causes(err) {
			sb.WriteRune('\n')
			printPretty(sb, depth+1, cause, opts)
		}
	}
}
