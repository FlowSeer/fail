package fail

import "time"

// PrinterOptions configures the behavior of a Printer.
//
// These options control which error fields are included in the output, formatting
// details such as indentation and color, and which metadata to display. This struct
// is typically used with functional options to customize printers.
type PrinterOptions struct {
	// Indent specifies the number of spaces to use for indentation of nested errors.
	Indent int
	// Color enables ANSI color output if true.
	// Printers may ignore this value if they do not support ANSI color output.
	Color bool
	// Time enables printing the error's timestamp if true.
	Time bool
	// TimeFormat specifies the layout for formatting the time, if Time is true.
	// Printers may ignore this value if they do not support custom time formats.
	TimeFormat string
	// Associated enables printing associated (non-causal) errors if true.
	Associated bool
	// Causes enables printing direct causes of the error if true.
	Causes bool
	// CauseDepth is the maximum recursion depth to print causes.
	// If 0, all causes are printed.
	CauseDepth int
	// Tags enables printing error tags if true.
	Tags bool
	// Attributes enables printing error attributes if true.
	Attributes bool
	// Code enables printing the error code if true.
	Code bool
	// Domain enables printing the error domain if true.
	Domain bool
	// ExitCode enables printing the process exit code if true.
	ExitCode bool
	// HttpStatusCode enables printing the HTTP status code if true.
	HttpStatusCode bool
	// UserMsg enables printing the user-facing message if true.
	UserMsg bool
	// TraceId enables printing the trace ID if true.
	TraceId bool
	// SpanId enables printing the span ID if true.
	SpanId bool
}

// DefaultOptions returns a PrinterOptions struct with all fields set to their default values.
//
// The defaults are suitable for most use cases, enabling all fields and using
// a standard indentation and time format.
func DefaultOptions() PrinterOptions {
	return PrinterOptions{
		Indent:         2,
		Color:          true,
		Time:           true,
		TimeFormat:     time.RFC3339,
		Associated:     true,
		Causes:         true,
		Tags:           true,
		Attributes:     true,
		Code:           true,
		Domain:         true,
		ExitCode:       true,
		HttpStatusCode: true,
		UserMsg:        true,
		TraceId:        true,
		SpanId:         true,
	}
}

// PrinterOption is a functional option for configuring PrinterOptions.
//
// Use PrinterOption functions to set fields on PrinterOptions when constructing
// or customizing a Printer.
type PrinterOption func(*PrinterOptions)

// PrintIndent sets the indentation level (number of spaces) for nested errors.
//
// Example: print.PrintIndent(4)
func PrintIndent(indent int) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Indent = indent
	}
}

// PrintColor enables or disables ANSI color output.
//
// Example: print.WithColor(false)
func PrintColor(color bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Color = color
	}
}

// PrintTime enables or disables printing the error's timestamp.
//
// Example: print.WithTime(false)
func PrintTime(time bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Time = time
	}
}

// PrintTimeFormat sets the layout for formatting the time, if time printing is enabled.
//
// Example: print.PrintTimeFormat(time.RFC1123)
func PrintTimeFormat(timeFormat string) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.TimeFormat = timeFormat
	}
}

// PrintAssociated enables or disables printing associated (non-causal) errors.
//
// Example: print.PrintAssociated(false)
func PrintAssociated(associated bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Associated = associated
	}
}

// PrintCauses enables or disables printing direct causes of the error.
//
// Example: print.PrintCauses(false)
func PrintCauses(causes bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Causes = causes
	}
}

// PrintCauseDepth sets the recursion depth of causes to print.
//
// Example: print.PrintTags(false)
func PrintCauseDepth(depth int) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.CauseDepth = depth
	}
}

// PrintTags enables or disables printing error tags.
//
// Example: print.PrintTags(false)
func PrintTags(tags bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Tags = tags
	}
}

// PrintAttributes enables or disables printing error attributes.
//
// Example: print.PrintAttributes(false)
func PrintAttributes(attributes bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Attributes = attributes
	}
}

// PrintCode enables or disables printing the error code.
//
// Example: print.PrintCode(false)
func PrintCode(code bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Code = code
	}
}

// PrintDomain enables or disables printing the error domain.
//
// Example: print.PrintDomain(false)
func PrintDomain(domain bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Domain = domain
	}
}

// PrintExitCode enables or disables printing the process exit code.
//
// Example: print.PrintExitCode(false)
func PrintExitCode(exitCode bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.ExitCode = exitCode
	}
}

// PrintHttpStatusCode enables or disables printing the HTTP status code.
//
// Example: print.PrintHttpStatusCode(false)
func PrintHttpStatusCode(httpStatusCode bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.HttpStatusCode = httpStatusCode
	}
}

// PrintUserMsg enables or disables printing the user-facing message.
//
// Example: print.PrintUserMsg(false)
func PrintUserMsg(userMsg bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.UserMsg = userMsg
	}
}

// PrintTraceId enables or disables printing the trace ID.
//
// Example: print.PrintTraceId(false)
func PrintTraceId(traceId bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.TraceId = traceId
	}
}

// PrintSpanId enables or disables printing the span ID.
//
// Example: print.PrintSpanId(false)
func PrintSpanId(spanId bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.SpanId = spanId
	}
}
