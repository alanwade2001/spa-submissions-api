package main

import "io"

// SubmissionService s
type SubmissionService struct {
	repositoryAPI RepositoryAPI
	initiationAPI    InitiationAPI
	validatorAPI  ValidatorAPI
}

// NewSubmissionService f
func NewSubmissionService(repositoryAPI RepositoryAPI, initiationAPI InitiationAPI, validatorAPI ValidatorAPI) SubmissionServiceAPI {
	service := SubmissionService{
		repositoryAPI: repositoryAPI,
		initiationAPI:    initiationAPI,
		validatorAPI:  validatorAPI,
	}
	return service
}

// CreateSubmission f
func (s SubmissionService) CreateSubmission(rc io.ReadCloser) (*Submission, error) {
	var initiation *Initiation
	var result *Result

	if data, err := io.ReadAll(rc); err != nil {
		return nil, err
	} else if initiation, err = s.initiationAPI.Parse(data); err != nil {
		return nil, err
	} else if result, err = s.validatorAPI.Validate(*initiation); err != nil {
		return nil, err
	}

	submission := Submission{
		Initiation:          *initiation,
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
