package fail

import (
	"encoding/json"
	"strings"
	"time"
)

// PrintJson returns a JSON-formatted string representation of the provided error.
//
// If the error is nil, this function returns the string "null" (the JSON null value).
//
// This function uses the default JSON printer with the given PrinterOptions to
// serialize the error and its metadata (such as causes, associated errors, codes, tags, etc.)
// into a JSON string. It is suitable for logging, diagnostics, or API responses.
//
// Example:
//
//	err := fail.New().Msg("something went wrong")
//	jsonStr := print.PrintJson(err)
//
// The output format and included fields can be customized using PrinterOptions.
func PrintJson(err error, opts ...PrinterOption) string {
	return JsonPrinter(opts...).Print(err)
}

// JsonPrinter returns a Printer that formats errors as JSON strings.
//
// The returned Printer serializes errors and their metadata into JSON, using the
// provided PrinterOptions to control which fields are included. If the error is nil,
// the Printer returns the string "null" (the JSON null value). This is useful for
// structured logging, diagnostics, or API error responses.
//
// Example:
//
//	printer := print.JsonPrinter(print.WithoutColor())
//	out := printer.Print(err)
func JsonPrinter(opts ...PrinterOption) Printer {
	return PrinterFunc(func(err error) string {
		return printJson(err, opts...)
	})
}

// printJson serializes the provided error into a JSON string according to the given PrinterOptions.
//
// This is an internal helper used by JsonPrinter and PrintJson. It panics if not implemented.
func printJson(err error, opts ...PrinterOption) string {
	if err == nil {
		return "null"
	}

	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	data := map[string]any{
		"msg": Message(err),
	}

	if o.Time {
		t := Time(err)
		if !t.IsZero() {
			timeFormat := time.RFC3339
			if o.TimeFormat != "" {
				timeFormat = o.TimeFormat
			}

			data["time"] = t.Format(timeFormat)
		}
	}

	if o.Associated {
		associated := Associated(err)
		if len(associated) > 0 {
			data["associated"] = associated
		}
	}

	if o.Causes {
		causes := Causes(err)
		if len(causes) > 0 {
			data["causes"] = causes
		}
	}

	if o.Tags {
		tags := Tags(err)
		if len(tags) > 0 {
			data["tags"] = tags
		}
	}

	if o.Attributes {
		attributes := Attributes(err)
		if len(attributes) > 0 {
			data["attributes"] = attributes
		}
	}

	if o.Code {
		code := Code(err)
		if code != "" {
			data["code"] = code
		}
	}

	if o.Domain {
		domain := Domain(err)
		if domain != "" {
			data["domain"] = domain
		}
	}

	if o.ExitCode {
		exitCode := ExitCode(err)
		if exitCode > 0 {
			data["exit_code"] = exitCode
		}
	}

	if o.HttpStatusCode {
		httpStatusCode := HttpStatusCode(err)
		if httpStatusCode > 0 {
			data["http_status_code"] = httpStatusCode
		}
	}

	if o.UserMsg {
		userMsg := UserMessage(err)
		if userMsg != "" {
			data["user_msg"] = userMsg
		}
	}

	if o.TraceId {
		traceId := TraceId(err)
		if traceId != "" {
			data["trace_id"] = traceId
		}
	}

	if o.SpanId {
		spanId := SpanId(err)
		if spanId != "" {
			data["span_id"] = spanId
		}
	}

	b, err := json.MarshalIndent(data, "", strings.Repeat(" ", o.Indent))
	if err != nil {
		panic(err)
	}

	return string(b)
}
