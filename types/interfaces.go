package types

import (
	"io"

	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"

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
	CreateSubmission(c *submission.SubmissionModel) (*submission.SubmissionModel, error)
	GetSubmission(id string) (*submission.SubmissionModel, error)
	GetSubmissions() ([]*submission.SubmissionModel, error)
}

// MessageAPI i
type MessageAPI interface {
	SendInitiation(i initiation.InitiationModel) error
}

// ConfigAPI si
type ConfigAPI interface {
	Load(path string) error
}

// XMLParserAPI i
type XMLParserAPI interface {
	Parse(data []byte) (types.Document, error)
}

// ValidatorAPI i
type ValidatorAPI interface {
	Validate(initiation initiation.InitiationModel) (*Result, error)
}

// InitiationMapperAPI i
type InitiationMapperAPI interface {
	Map(types.Document) (*initiation.InitiationModel, error)
}

// GroupHeaderMapperAPI i
type GroupHeaderMapperAPI interface {
	Map(*xpath.Context) (*initiation.GroupHeaderReference, error)
}

// PaymentInformationMapperAPI i
type PaymentInformationMapperAPI interface {
	Map(*xpath.Context) ([]*initiation.PaymentInstructionReference, error)
}

// SubmissionServiceAPI i
type SubmissionServiceAPI interface {
	CreateSubmission(io.ReadCloser, submission.UserReference) (*submission.SubmissionModel, error)
	GetSubmission(ID string) (*submission.SubmissionModel, error)
	GetSubmissions() ([]*submission.SubmissionModel, error)
}

// InitiationAPI i
type InitiationAPI interface {
	Parse(data []byte) (*initiation.InitiationModel, error)
}

// CustomerAPI i
type CustomerAPI interface {
	Find(user submission.UserReference) (*submission.CustomerReference, error)
}

// UserAPI i
type UserAPI interface {
	Find(*gin.Context) (*submission.UserReference, error)
}
