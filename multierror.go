package multierror

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Errors      []error
	ErrorFormat ErrorFormatFunc
}

func (e *Error) Error() string {
	fn := e.ErrorFormat
	if fn == nil {
		fn = ListFormatFunc
	}

	return fn(e.Errors)
}

// MarshalJSON returns a valid json representation of a multierror,
// as an object with an array of error strings.
func (e *Error) MarshalJSON() ([]byte, error) {
	j := map[string][]string{
		"errors": []string{},
	}
	for _, err := range e.Errors {
		j["errors"] = append(j["errors"], err.Error())
	}

	return json.Marshal(j)
}

// UnmarshalJSON parses the output of Marshal json.
func (e *Error) UnmarshalJSON(b []byte) error {
	j := make(map[string][]string)
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	if _, ok := j["errors"]; ok {
		for _, msg := range j["errors"] {
			e.Errors = append(e.Errors, errors.New(msg))
		}
	}
	return nil
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *Error) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if len(e.Errors) == 0 {
		return nil
	}

	return e
}

func (e *Error) GoString() string {
	return fmt.Sprintf("*%#v", *e)
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the errwrap.Wrapper interface so that
// multierror.Error can be used with that library.
//
// This method is not safe to be called concurrently and is no different
// than accessing the Errors field directly. It is implemented only to
// satisfy the errwrap.Wrapper interface.
func (e *Error) WrappedErrors() []error {
	return e.Errors
}
