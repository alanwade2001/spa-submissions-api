package main

import types "github.com/alanwade2001/spa-common"

// Submission s
type Submission struct {
	ID               string `bson:"_id"`
	Initiation       *types.Initiation
	ValidationResult Result
	Customer         *types.CustomerReference
	Submitter        types.UserReference
}

// Submissions a
type Submissions []Submission

// Failure s
type Failure struct {
	Error error
	Level FailureLevel
}

// Failures a
type Failures []Failure

// FailureLevel i
type FailureLevel int

const (
	// LevelCustomer c
	LevelCustomer FailureLevel = iota
	// LevelGroupHeader c
	LevelGroupHeader FailureLevel = iota
	// LevelPayment c
	LevelPayment FailureLevel = iota
	// LevelTransaction c
	LevelTransaction FailureLevel = iota
)

// Result s
type Result struct {
	Success  bool
	Failures Failures
}
