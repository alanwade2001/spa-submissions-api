package main

// Submission s
type Submission struct {
	ID               string `bson:"_id"`
	Initiation       *Initiation
	ValidationResult Result
	Customer         *Customer
	Submitter        User
}

// Submissions a
type Submissions []Submission

// User s
type User struct {
	Email string
}

// Users a
type Users []User

// Roles s
type Roles struct {
	Submitters Users
	Approvers  Users
	Admins     Users
}

// Customer s
type Customer struct {
	ID                string
	Name              string
	InitiatingPartyID string
}

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

// Initiation s
type Initiation struct {
	GroupHeader         GroupHeader
	PaymentInstructions PaymentInstructions
	Customer            *Customer
}

// GroupHeader s
type GroupHeader struct {
	MessageID            string
	CreationDateTime     string
	NumberOfTransactions string
	ControlSum           string
	InitiatingPartyID    string
}

// PaymentInstruction s
type PaymentInstruction struct {
	PaymentID              string
	NumberOfTransactions   string
	ControlSum             string
	RequestedExecutionDate string
	Debtor                 Account
}

// PaymentInstructions a
type PaymentInstructions []PaymentInstruction

// Account s
type Account struct {
	Name string
	IBAN string
	BIC  string
}
