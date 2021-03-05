package main

import "io"

// SubmissionService s
type SubmissionService struct {
	repositoryAPI RepositoryAPI
	pain001API    Pain001API
	validatorAPI  ValidatorAPI
}

// NewSubmissionService f
func NewSubmissionService(repositoryAPI RepositoryAPI, pain001API Pain001API, validatorAPI ValidatorAPI) SubmissionServiceAPI {
	service := SubmissionService{
		repositoryAPI: repositoryAPI,
		pain001API:    pain001API,
		validatorAPI:  validatorAPI,
	}
	return service
}

// CreateSubmission f
func (s SubmissionService) CreateSubmission(rc io.ReadCloser) (*Submission, error) {
	var pain001 *Pain001
	var result *Result

	if data, err := io.ReadAll(rc); err != nil {
		return nil, err
	} else if pain001, err = s.pain001API.Parse(data); err != nil {
		return nil, err
	} else if result, err = s.validatorAPI.Validate(*pain001); err != nil {
		return nil, err
	}

	submission := Submission{
		Pain001:          *pain001,
		ValidationResult: *result,
	}

	return &submission, nil
}

// GetSubmission f
func (s SubmissionService) GetSubmission(ID string) (*Submission, error) {
	return s.repositoryAPI.GetSubmission(ID)
}

// GetSubmissions f
func (s SubmissionService) GetSubmissions() (*Submissions, error) {
	return s.repositoryAPI.GetSubmissions()
}
