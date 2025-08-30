package printer

// Printer is an interface for formatting errors as strings.
//
// Implementations of Printer can be used to customize how errors are rendered
// for logs, user interfaces, or diagnostics. The Print method should return a
// string representation of the provided error, potentially including details
// such as causes, associated errors, codes, tags, and more.
type Printer interface {
	Print(err error) string
}

// PrinterFunc is an adapter to allow the use of ordinary functions as Printers.
//
// Any function with the appropriate signature can be converted to a Printer
// by using PrinterFunc(f). This enables flexible and concise custom printers.
type PrinterFunc func(err error) string

// Print calls the underlying function to print the error.
func (f PrinterFunc) Print(err error) string {
	return f(err)
}
