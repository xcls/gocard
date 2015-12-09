package valid

import "fmt"

type validator struct {
	errors []error
}

func NewValidator() *validator {
	return &validator{[]error{}}
}

func (v *validator) Errors() []error {
	return v.errors
}

func (v *validator) AddError(err error) {
	v.errors = append(v.errors, err)
}

func (v *validator) ValidateMinLength(label, val string, length int) {
	if len(val) >= length {
		return
	}
	if length <= 1 {
		v.AddError(fmt.Errorf("%s can't be blank", label))
	} else {
		v.AddError(fmt.Errorf("%s must be at least %d characters long", label, length))
	}
}

func (v *validator) ValidateConfirmation(label, original string, confirmation string) {
	if original != confirmation {
		v.AddError(fmt.Errorf("%s and its confirmation don't match", label))
	}
}
