# go-multierror

[![CircleCI](https://img.shields.io/circleci/build/github/hashicorp/go-multierror/master)](https://circleci.com/gh/hashicorp/go-multierror)
[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/go-multierror.svg)](https://pkg.go.dev/github.com/hashicorp/go-multierror)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hashicorp/go-multierror)

[circleci]: https://app.circleci.com/pipelines/github/hashicorp/go-multierror
[godocs]: https://pkg.go.dev/github.com/hashicorp/go-multierror

> **Note**: As of Go 1.20, the standard library provides [`errors.Join`](https://pkg.go.dev/errors#Join) which offers similar functionality for combining multiple errors. **For new projects, we recommend using `errors.Join` from the standard library.** This package provides additional features like custom formatting, the `Group` pattern for concurrent error collection, and utility functions (`Append`, `Flatten`, `Prefix`).

## Migrating to `errors.Join`

**Basic error aggregation:**
```go
// Before (go-multierror)
var result error
result = multierror.Append(result, err1)
result = multierror.Append(result, err2)
return result

// After (stdlib)
var errs []error
if err1 != nil {
    errs = append(errs, err1)
}
if err2 != nil {
    errs = append(errs, err2)
}
if len(errs) > 0 {
    return errors.Join(errs...)
}
return nil
```

**Custom formatting:**
```go
// go-multierror allows custom formatting
// With stdlib, wrap the joined error with a custom type if needed
type FormattedError struct {
    errs []error
}

func (e *FormattedError) Error() string {
    // Your custom format here
    return fmt.Sprintf("multiple errors: %v", e.errs)
}

func (e *FormattedError) Unwrap() []error {
    return e.errs
}
```

**Concurrent error collection (Group):**
```go
// golang.org/x/sync/errgroup returns only the first error:
import "golang.org/x/sync/errgroup"

g := new(errgroup.Group)
g.Go(func() error { return task1() })
g.Go(func() error { return task2() })
err := g.Wait() // Returns first error only

// To collect ALL errors like multierror.Group does, use a mutex:
type ErrorCollector struct {
    mu   sync.Mutex
    errs []error
}

func (ec *ErrorCollector) Add(err error) {
    if err != nil {
        ec.mu.Lock()
        ec.errs = append(ec.errs, err)
        ec.mu.Unlock()
    }
}

func (ec *ErrorCollector) Err() error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    if len(ec.errs) == 0 {
        return nil
    }
    return errors.Join(ec.errs...)
}
```

---

`go-multierror` is a package for Go that provides a mechanism for
representing a list of `error` values as a single `error`.

This allows a function in Go to return an `error` that might actually
be a list of errors. If the caller knows this, they can unwrap the
list and access the errors. If the caller doesn't know, the error
formats to a nice human-readable format.

`go-multierror` is compatible with the Go standard library
[errors](https://golang.org/pkg/errors/) package, supporting the
`As`, `Is`, and `Unwrap` functions for error introspection. Note that `Unwrap` 
returns errors one at a time via chaining, unlike `errors.Join` which 
implements the newer `Unwrap() []error` signature (Go 1.20+).

## Installation and Docs

Install using `go get github.com/hashicorp/go-multierror`.

Full documentation is available at
https://pkg.go.dev/github.com/hashicorp/go-multierror

### Requires go version 1.13 or newer

`go-multierror` requires go version 1.13 or newer. Go 1.13 introduced
[error wrapping](https://golang.org/doc/go1.13#error_wrapping), which
this library takes advantage of.

If you need to use an earlier version of go, you can use the
[v1.0.0](https://github.com/hashicorp/go-multierror/tree/v1.0.0)
tag, which doesn't rely on features in go 1.13.

If you see compile errors that look like the below, it's likely that
you're on an older version of go:

```
/go/src/github.com/hashicorp/go-multierror/multierror.go:112:9: undefined: errors.As
/go/src/github.com/hashicorp/go-multierror/multierror.go:117:9: undefined: errors.Is
```

## Usage

go-multierror is easy to use and purposely built to be unobtrusive in
existing Go applications/libraries that may not be aware of it.

**Building a list of errors**

The `Append` function is used to create a list of errors. This function
behaves a lot like the Go built-in `append` function: it doesn't matter
if the first argument is nil, a `multierror.Error`, or any other `error`,
the function behaves as you would expect.

```go
var result error

if err := step1(); err != nil {
	result = multierror.Append(result, err)
}
if err := step2(); err != nil {
	result = multierror.Append(result, err)
}

return result
```

**Customizing the formatting of the errors**

By specifying a custom `ErrorFormat`, you can customize the format
of the `Error() string` function:

```go
var result *multierror.Error

// ... accumulate errors here, maybe using Append

if result != nil {
	result.ErrorFormat = func([]error) string {
		return "errors!"
	}
}
```

**Accessing the list of errors**

`multierror.Error` implements `error` so if the caller doesn't know about
multierror, it will work just fine. But if you're aware a multierror might
be returned, you can use type switches to access the list of errors:

```go
if err := something(); err != nil {
	if merr, ok := err.(*multierror.Error); ok {
		// Use merr.Errors
	}
}
```

You can also use the standard [`errors.Unwrap`](https://golang.org/pkg/errors/#Unwrap)
function. This will continue to unwrap into subsequent errors until none exist.

**Extracting an error**

The standard library [`errors.As`](https://golang.org/pkg/errors/#As)
function can be used directly with a multierror to extract a specific error:

```go
// Assume err is a multierror value
err := somefunc()

// We want to know if "err" has a "RichErrorType" in it and extract it.
var errRich RichErrorType
if errors.As(err, &errRich) {
	// It has it, and now errRich is populated.
}
```

**Checking for an exact error value**

Some errors are returned as exact errors such as the [`ErrNotExist`](https://golang.org/pkg/os/#pkg-variables)
error in the `os` package. You can check if this error is present by using
the standard [`errors.Is`](https://golang.org/pkg/errors/#Is) function.

```go
// Assume err is a multierror value
err := somefunc()
if errors.Is(err, os.ErrNotExist) {
	// err contains os.ErrNotExist
}
```

**Returning a multierror only if there are errors**

If you build a `multierror.Error`, you can use the `ErrorOrNil` function
to return an `error` implementation only if there are errors to return:

```go
var result *multierror.Error

// ... accumulate errors here

// Return the `error` only if errors were added to the multierror, otherwise
// return nil since there are no errors.
return result.ErrorOrNil()
```
