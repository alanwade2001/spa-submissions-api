package types

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
