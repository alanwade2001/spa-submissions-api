package main

import "io"

// SubmissionService s
type SubmissionService struct {
	repositoryAPI RepositoryAPI
	initiationAPI InitiationAPI
	validatorAPI  ValidatorAPI
	customerAPI   CustomerAPI
}

// NewSubmissionService f
func NewSubmissionService(repositoryAPI RepositoryAPI, initiationAPI InitiationAPI, validatorAPI ValidatorAPI, customerAPI CustomerAPI) SubmissionServiceAPI {
	service := SubmissionService{
		repositoryAPI: repositoryAPI,
		initiationAPI: initiationAPI,
		validatorAPI:  validatorAPI,
		customerAPI:   customerAPI,
	}
	return service
}

// CreateSubmission f
func (s SubmissionService) CreateSubmission(rc io.ReadCloser, user User) (submission *Submission, err error) {
	var initiation *Initiation
	var result *Result
	var data []byte
	var customer *Customer

	if customer, err = s.customerAPI.Find(user); err != nil {
		return nil, err
	} else if data, err = io.ReadAll(rc); err != nil {
		return nil, err
	} else if initiation, err = s.initiationAPI.Parse(data); err != nil {
		return nil, err
	} else {
		initiation.Customer = customer
	}

	if result, err = s.validatorAPI.Validate(*initiation); err != nil {
		return nil, err
	}

	submission = &Submission{
		Initiation:       initiation,
		ValidationResult: *result,
		Customer:         customer,
		Submitter:        user,
	}

	if submission, err = s.repositoryAPI.CreateSubmission(submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// GetSubmission f
func (s SubmissionService) GetSubmission(ID string) (*Submission, error) {
	return s.repositoryAPI.GetSubmission(ID)
}

// GetSubmissions f
func (s SubmissionService) GetSubmissions() (*Submissions, error) {
	return s.repositoryAPI.GetSubmissions()
}
