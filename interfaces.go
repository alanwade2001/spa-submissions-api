package main

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
)

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

// XMLParserAPI i
type XMLParserAPI interface {
	Parse(data []byte) (types.Document, error)
}

// ValidatorAPI i
type ValidatorAPI interface {
	Validate(pain001 Pain001) (*Result, error)
}

// Pain001MapperAPI i
type Pain001MapperAPI interface {
	Map(types.Document) (*Pain001, error)
}

// GroupHeaderMapperAPI i
type GroupHeaderMapperAPI interface {
	Map(*xpath.Context) (*GroupHeader, error)
}

// PaymentInformationMapperAPI i
type PaymentInformationMapperAPI interface {
	Map(*xpath.Context) (PmtInfs, error)
}

// SubmissionServiceAPI i
type SubmissionServiceAPI interface {
	CreateSubmission(io.ReadCloser) (*Submission, error)
	GetSubmission(ID string) (*Submission, error)
	GetSubmissions() (*Submissions, error)
}

// Pain001API i
type Pain001API interface {
	Parse(data []byte) (*Pain001, error)
}
