package main

import types "github.com/alanwade2001/spa-common"

// Validator s
type Validator struct {
}

// NewValidator f
func NewValidator() ValidatorAPI {
	return Validator{}
}

// Validate f
func (v Validator) Validate(initiation types.Initiation) (*Result, error) {
	result := Result{Success: true}

	return &result, nil
}
