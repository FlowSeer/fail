package printer

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/FlowSeer/fail"
)

// Json returns a JSON-formatted string representation of the provided error.
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
//	jsonStr := print.Json(err)
//
// The output format and included fields can be customized using PrinterOptions.
func Json(err error, opts ...PrinterOption) string {
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
// This is an internal helper used by JsonPrinter and Json. It panics if not implemented.
func printJson(err error, opts ...PrinterOption) string {
	if err == nil {
		return "null"
	}

	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	data := map[string]any{
		"msg": fail.Message(err),
	}

	if o.Time {
		t := fail.Time(err)
		if !t.IsZero() {
			timeFormat := time.RFC3339
			if o.TimeFormat != "" {
				timeFormat = o.TimeFormat
			}

			data["time"] = t.Format(timeFormat)
		}
	}

	if o.Associated {
		associated := fail.Associated(err)
		if len(associated) > 0 {
			data["associated"] = associated
		}
	}

	if o.Causes {
		causes := fail.Causes(err)
		if len(causes) > 0 {
			data["causes"] = causes
		}
	}

	if o.Tags {
		tags := fail.Tags(err)
		if len(tags) > 0 {
			data["tags"] = tags
		}
	}

	if o.Attributes {
		attributes := fail.Attributes(err)
		if len(attributes) > 0 {
			data["attributes"] = attributes
		}
	}

	if o.Code {
		code := fail.Code(err)
		if code != "" {
			data["code"] = code
		}
	}

	if o.Domain {
		domain := fail.Domain(err)
		if domain != "" {
			data["domain"] = domain
		}
	}

	if o.ExitCode {
		exitCode := fail.ExitCode(err)
		if exitCode > 0 {
			data["exit_code"] = exitCode
		}
	}

	if o.HttpStatusCode {
		httpStatusCode := fail.HttpStatusCode(err)
		if httpStatusCode > 0 {
			data["http_status_code"] = httpStatusCode
		}
	}

	if o.UserMsg {
		userMsg := fail.UserMessage(err)
		if userMsg != "" {
			data["user_msg"] = userMsg
		}
	}

	if o.TraceId {
		traceId := fail.TraceId(err)
		if traceId != "" {
			data["trace_id"] = traceId
		}
	}

	if o.SpanId {
		spanId := fail.SpanId(err)
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
