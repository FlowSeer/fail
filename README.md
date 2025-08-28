# FlowSeer/fail

A rich, fluent error handling library for Go that provides comprehensive error context, tracing, and formatting capabilities.

## Features

- **Rich Error Context**: Add codes, exit codes, HTTP status codes, domains, tags, and attributes
- **Error Relationships**: Track causes and associated errors
- **User Messages**: Separate developer and user-facing error messages
- **Tracing Support**: Built-in OpenTelemetry trace and span ID support
- **Fluent Builder API**: Chain method calls for easy error construction
- **Multiple Output Formats**: JSON and pretty-printed error representations
- **Domain Categorization**: Predefined domains for common error types
- **Shortcut Functions**: Quick error creation for common use cases

## Installation

```bash
go get github.com/FlowSeer/fail
```

## Quick Start

### Basic Error Creation

```go
package main

import "github.com/FlowSeer/fail"

func main() {
    // Simple error
    err := fail.Msg("database connection failed")
    
    // Formatted error
    err = fail.Msgf("failed to connect to %s on port %d", "localhost", 5432)
    
    // Wrap existing error
    err = fail.Wrap("failed to read file", io.EOF)
}
```

### Rich Error with Context

```go
err := fail.New().
    UserMsg("Unable to process your request. Please try again.").
    Code("DB_CONNECTION_ERROR").
    Domain(fail.DomainDatabase).
    ExitCode(1).
    HttpStatusCode(503).
    Tag(fail.TagNetwork, fail.TagTimeout).
    Attribute("host", "db.example.com").
    Attribute("port", 5432).
    Cause(io.EOF).
    TraceId("abcdef1234567890").
    SpanId("1234567890abcdef").
    Msg("database connection failed")
```

## API Reference

### Builder Methods

The `fail.New()` builder provides a fluent interface:

```go
err := fail.New().
    Time(time.Now()).                    // Set timestamp
    UserMsg("User-friendly message").    // Set user message
    Code("ERROR_CODE").                  // Set error code
    Domain(fail.DomainDatabase).         // Set domain
    ExitCode(1).                         // Set exit code
    HttpStatusCode(500).                 // Set HTTP status code
    Tag("network", "timeout").           // Add tags
    Attribute("key", "value").           // Add attributes
    Cause(originalError).                // Add cause
    Associate(relatedError).             // Add associated error
    TraceId("trace-id").                 // Set trace ID
    SpanId("span-id").                   // Set span ID
    Msg("Developer message")             // Set message and build
```

### Shortcut Functions

Quick error creation:

```go
// Simple errors
err := fail.Msg("something went wrong")
err := fail.Msgf("failed to %s", "connect")

// Wrapping errors
err := fail.Wrap(originalError, "operation failed")
err := fail.Wrapf(originalError, "failed to %s", "process")

// Multiple causes
err := fail.WrapMany("multiple errors", err1, err2, err3)
```

### Error Inspection

Extract information from errors:

```go
// Basic information
msg := fail.Message(err)
userMsg := fail.UserMessage(err)
code := fail.Code(err)
domain := fail.Domain(err)

// Status codes
exitCode := fail.ExitCode(err)
httpStatus := fail.HttpStatusCode(err)

// Relationships
causes := fail.Causes(err)
associated := fail.Associated(err)

// Metadata
tags := fail.Tags(err)
attrs := fail.Attributes(err)
traceId := fail.TraceId(err)
spanId := fail.SpanId(err)
```

## Error Printing

### Pretty Printing

```go
import "github.com/FlowSeer/fail/print"

// Pretty print with default options
output := print.Pretty(err)
fmt.Println(output)

// Custom pretty printing
output := print.PrettyWithOptions(err, print.Options{
    ShowCauses:     true,
    ShowAssociated: true,
    ShowTags:       true,
    ShowAttributes: true,
    ShowTrace:      true,
})
```

### JSON Output

```go
import "github.com/FlowSeer/fail/print"

// JSON output
output := print.JSON(err)
fmt.Println(output)

// Pretty JSON
output := print.JSONPretty(err)
fmt.Println(output)
```

## Example

```go
// Simple error
err := fail.Msg("database connection failed")

// Rich error with context
err := fail.New().
    UserMsg("Unable to process your request. Please try again.").
    Code("DB_CONNECTION_ERROR").
    Domain(fail.DomainDatabase).
    HttpStatusCode(503).
    Tag(fail.TagNetwork, fail.TagTimeout).
    Cause(originalError).
    Msg("database connection failed")

// Wrap existing error
err := fail.Wrap("failed to read file", io.EOF)
```