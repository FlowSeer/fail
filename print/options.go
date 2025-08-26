package print

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

// WithIndent sets the indentation level (number of spaces) for nested errors.
//
// Example: print.WithIndent(4)
func WithIndent(indent int) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Indent = indent
	}
}

// WithColor enables or disables ANSI color output.
//
// Example: print.WithColor(false)
func WithColor(color bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Color = color
	}
}

// WithTime enables or disables printing the error's timestamp.
//
// Example: print.WithTime(false)
func WithTime(time bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Time = time
	}
}

// WithTimeFormat sets the layout for formatting the time, if time printing is enabled.
//
// Example: print.WithTimeFormat(time.RFC1123)
func WithTimeFormat(timeFormat string) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.TimeFormat = timeFormat
	}
}

// WithAssociated enables or disables printing associated (non-causal) errors.
//
// Example: print.WithAssociated(false)
func WithAssociated(associated bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Associated = associated
	}
}

// WithCauses enables or disables printing direct causes of the error.
//
// Example: print.WithCauses(false)
func WithCauses(causes bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Causes = causes
	}
}

// WithTags enables or disables printing error tags.
//
// Example: print.WithTags(false)
func WithTags(tags bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Tags = tags
	}
}

// WithAttributes enables or disables printing error attributes.
//
// Example: print.WithAttributes(false)
func WithAttributes(attributes bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Attributes = attributes
	}
}

// WithCode enables or disables printing the error code.
//
// Example: print.WithCode(false)
func WithCode(code bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Code = code
	}
}

// WithDomain enables or disables printing the error domain.
//
// Example: print.WithDomain(false)
func WithDomain(domain bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.Domain = domain
	}
}

// WithExitCode enables or disables printing the process exit code.
//
// Example: print.WithExitCode(false)
func WithExitCode(exitCode bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.ExitCode = exitCode
	}
}

// WithHttpStatusCode enables or disables printing the HTTP status code.
//
// Example: print.WithHttpStatusCode(false)
func WithHttpStatusCode(httpStatusCode bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.HttpStatusCode = httpStatusCode
	}
}

// WithUserMsg enables or disables printing the user-facing message.
//
// Example: print.WithUserMsg(false)
func WithUserMsg(userMsg bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.UserMsg = userMsg
	}
}

// WithTraceId enables or disables printing the trace ID.
//
// Example: print.WithTraceId(false)
func WithTraceId(traceId bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.TraceId = traceId
	}
}

// WithSpanId enables or disables printing the span ID.
//
// Example: print.WithSpanId(false)
func WithSpanId(spanId bool) PrinterOption {
	return func(opts *PrinterOptions) {
		opts.SpanId = spanId
	}
}
