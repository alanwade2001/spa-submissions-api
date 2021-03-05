package main

// Validator s
type Validator struct {
}

// NewValidator f
func NewValidator() ValidatorAPI {
	return Validator{}
}

// Validate f
func (v Validator) Validate(pain001 Pain001) (*Result, error) {
	result := Result{Success: true}
	
	return &result, nil
}
