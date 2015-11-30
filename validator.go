package main

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

func (v *validator) addError(err error) {
	v.errors = append(v.errors, err)
}

func (v *validator) ValidateMinLength(label, val string, length int) {
	if len(val) >= length {
		return
	}
	if length <= 1 {
		v.addError(fmt.Errorf("%s can't be blank", label))
	} else {
		v.addError(fmt.Errorf("%s must be at least %d characters long", label, length))
	}
}
