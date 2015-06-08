package multierror

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestFlatten(t *testing.T) {
	original := &Error{
		Errors: []error{
			errors.New("one"),
			&Error{
				Errors: []error{
					errors.New("two"),
					&Error{
						Errors: []error{
							errors.New("three"),
						},
					},
				},
			},
		},
	}

	expected := strings.TrimSpace(`
3 error(s) occurred:

* one
* two
* three
	`)
	actual := fmt.Sprintf("%s", Flatten(original))

	if expected != actual {
		t.Fatalf("expected: %s, got: %s", expected, actual)
	}
}
