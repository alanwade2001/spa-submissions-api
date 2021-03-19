package services

import (
	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/types"
)

// Validator s
type Validator struct {
}

// NewValidator f
func NewValidator() types.ValidatorAPI {
	return Validator{}
}

// Validate f
func (v Validator) Validate(initiation initiation.InitiationModel) (*types.Result, error) {
	result := types.Result{Success: true}

	return &result, nil
}
