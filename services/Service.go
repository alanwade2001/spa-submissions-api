package services

import (
	"io"
	"time"

	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
)

// SubmissionService s
type SubmissionService struct {
	repositoryAPI  types.RepositoryAPI
	initiationAPI  types.InitiationAPI
	validatorAPI   types.ValidatorAPI
	customerAPI    types.CustomerAPI
	messageAPI     types.MessageAPI
	idGeneratorAPI types.IdGeneratorAPI
}

// NewSubmissionService f
func NewSubmissionService(idGeneratorAPI types.IdGeneratorAPI, repositoryAPI types.RepositoryAPI, initiationAPI types.InitiationAPI, validatorAPI types.ValidatorAPI, customerAPI types.CustomerAPI, messageAPI types.MessageAPI) types.SubmissionServiceAPI {
	service := SubmissionService{
		idGeneratorAPI: idGeneratorAPI,
		repositoryAPI:  repositoryAPI,
		initiationAPI:  initiationAPI,
		validatorAPI:   validatorAPI,
		customerAPI:    customerAPI,
		messageAPI:     messageAPI,
	}
	return service
}

// CreateSubmission f
func (s SubmissionService) CreateSubmission(rc io.ReadCloser, user submission.UserReference) (model *submission.SubmissionModel, err error) {
	var init *initiation.InitiationModel
	//var result *types.Result
	var data []byte
	var customer *submission.CustomerReference

	if customer, err = s.customerAPI.Find(user); err != nil {
		return nil, err
	} else if data, err = io.ReadAll(rc); err != nil {
		return nil, err
	} else if init, err = s.initiationAPI.Parse(data); err != nil {
		return nil, err
	}

	if _, err = s.validatorAPI.Validate(*init); err != nil {
		return nil, err
	}

	model = &submission.SubmissionModel{
		GroupHeader: &submission.GroupHeaderReference{
			ControlSum:           init.GroupHeader.ControlSum,
			CreationDateTime:     init.GroupHeader.CreationDateTime,
			InitiatingPartyId:    init.GroupHeader.InitiatingPartyId,
			MessageId:            init.GroupHeader.MessageId,
			NumberOfTransactions: init.GroupHeader.NumberOfTransactions,
		},
		Customer:    customer,
		SubmittedBy: &user,
		SubmittedAt: time.Now().String(),
	}

	if model, err = s.repositoryAPI.CreateSubmission(model); err != nil {
		return nil, err
	}

	init.Id = s.idGeneratorAPI.Next()

	init.Customer = &initiation.CustomerReference{
		CustomerId:   customer.Id,
		CustomerName: customer.Name,
	}

	if err = s.messageAPI.SendInitiation(*init); err != nil {
		return nil, err
	}

	return model, nil
}

// GetSubmission f
func (s SubmissionService) GetSubmission(ID string) (*submission.SubmissionModel, error) {
	return s.repositoryAPI.GetSubmission(ID)
}

// GetSubmissions f
func (s SubmissionService) GetSubmissions() ([]*submission.SubmissionModel, error) {
	return s.repositoryAPI.GetSubmissions()
}
