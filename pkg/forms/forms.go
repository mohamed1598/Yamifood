package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	err := make(map[string][]string)
	return &Form{
		data,
		errors(err),
	}
}
func (f *Form) HasRequired(tagIds ...string) {
	for _, tagId := range tagIds {
		value := f.Get(tagId)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(tagId, "This filled can't be blank")
		}
	}
}

func (f *Form) MinLength(tagId string, length int, r *http.Request) bool {
	x := r.Form.Get(tagId)
	if len(x) < length {
		f.Errors.AddError(tagId, fmt.Sprintf("this field must be %d characters or more", length))
		return false
	}
	return true
}
func (f *Form) HasValue(tagId string, r *http.Request) bool {
	x := r.Form.Get(tagId)
	return x != ""
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
func (f *Form) IsEmail(tagId string) {
	if !govalidator.IsEmail(f.Get(tagId)) {
		f.Errors.AddError(tagId, "Invalid Email")
	}
}
