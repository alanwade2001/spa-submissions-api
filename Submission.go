package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmissionRouter s
type SubmissionRouter struct {
	serviceAPI SubmissionServiceAPI
}

// NewSubmissionRouter f
func NewSubmissionRouter(repositoryAPI RepositoryAPI, serviceAPI SubmissionServiceAPI) SubmissionAPI {

	submissionAPI := SubmissionRouter{serviceAPI}

	return &submissionAPI
}

// CreateSubmission f
func (cr *SubmissionRouter) CreateSubmission(ctx *gin.Context) {

	if submission, err := cr.serviceAPI.CreateSubmission(ctx.Request.Body); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusCreated, submission)
	}

}

// GetSubmission f
func (cr *SubmissionRouter) GetSubmission(ctx *gin.Context) {
	submissionID := ctx.Param("id")

	if submission, err := cr.serviceAPI.GetSubmission(submissionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, submission)
	}
}

// GetSubmissions f
func (cr *SubmissionRouter) GetSubmissions(ctx *gin.Context) {
	if submissions, err := cr.serviceAPI.GetSubmissions(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, submissions)
	}
}
