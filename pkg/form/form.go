package form

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func NewForm(values url.Values) *Form {
	return &Form{
		values,
		map[string][]string{},
	}
}

// check if the value of specified filed is empty
func (f *Form) Require(fields ...string) {
	for _, field := range fields {
		val := f.Values.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field cannot be black")
		}
	}

}

// check if the value size of specified field exceeds the limit
func (f *Form) MaxLength(field string, limit int) {
	val := f.Values.Get(field)

	if utf8.RuneCountInString(val) > limit {
		f.Errors.Add(field, fmt.Sprintf("This field is limited with %d words", limit))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	val := f.Values.Get(field)

	for _, opt := range opts {
		if val == opt {
			return
		}
	}

	f.Errors.Add(field, "This value is invalid")
}

// check if there is any error in the form data
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
