package main

// Submission s
type Submission struct {
	ID                string `bson:"_id"`
	Name              string
	InitiatingPartyID string
	Roles             Roles
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
