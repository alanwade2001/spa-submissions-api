package main

import "github.com/gin-gonic/gin"

// SubmissionAPI i
type SubmissionAPI interface {
	CreateSubmission(*gin.Context)
	GetSubmission(*gin.Context)
	GetSubmissions(*gin.Context)
}

// ServerAPI i
type ServerAPI interface {
	Run() error
}

// RegisterAPI i
type RegisterAPI interface {
	Register() error
}

// RepositoryAPI i
type RepositoryAPI interface {
	CreateSubmission(c *Submission) (*Submission, error)
	GetSubmission(id string) (*Submission, error)
	GetSubmissions() (*Submissions, error)
}

// ConfigAPI si
type ConfigAPI interface {
	Load() error
}
