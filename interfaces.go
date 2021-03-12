package main

import (
	"io"

	spatypes "github.com/alanwade2001/spa-common"
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

// MessageAPI i
type MessageAPI interface {
	SendInitiation(i spatypes.Initiation) error
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
	Validate(initiation spatypes.Initiation) (*Result, error)
}

// InitiationMapperAPI i
type InitiationMapperAPI interface {
	Map(types.Document) (*spatypes.Initiation, error)
}

// GroupHeaderMapperAPI i
type GroupHeaderMapperAPI interface {
	Map(*xpath.Context) (*spatypes.GroupHeader, error)
}

// PaymentInformationMapperAPI i
type PaymentInformationMapperAPI interface {
	Map(*xpath.Context) (*[]spatypes.PaymentInstruction, error)
}

// SubmissionServiceAPI i
type SubmissionServiceAPI interface {
	CreateSubmission(io.ReadCloser, spatypes.UserReference) (*Submission, error)
	GetSubmission(ID string) (*Submission, error)
	GetSubmissions() (*Submissions, error)
}

// InitiationAPI i
type InitiationAPI interface {
	Parse(data []byte) (*spatypes.Initiation, error)
}

// CustomerAPI i
type CustomerAPI interface {
	Find(user spatypes.UserReference) (*spatypes.CustomerReference, error)
}

// UserAPI i
type UserAPI interface {
	Find(*gin.Context) (*spatypes.UserReference, error)
}
