package main

// Submission s
type Submission struct {
	ID               string `bson:"_id"`
	Initiation       Initiation
	ValidationResult Result
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
	GroupHeader GroupHeader
	PmtInfs     PmtInfs
}

// GroupHeader s
type GroupHeader struct {
	MsgID    string
	CreDtTm  string
	NbOfTxs  string
	CtrlSum  string
	InitgPty string
}

// PmtInf s
type PmtInf struct {
	PmtInfID    string
	NbOfTxs     string
	CtrlSum     string
	ReqdExctnDt string
	Dbtr        Account
}

// PmtInfs a
type PmtInfs []PmtInf

// Account s
type Account struct {
	Name string
	IBAN string
	BIC  string
}
