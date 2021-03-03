package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmissionRouter s
type SubmissionRouter struct {
	repositoryAPI RepositoryAPI
}


// NewSubmissionRouter f
func NewSubmissionRouter(repositoryAPI RepositoryAPI) SubmissionAPI {

	submissionAPI := SubmissionRouter{repositoryAPI}

	return &submissionAPI
}

// CreateSubmission f
func (cr *SubmissionRouter) CreateSubmission(ctx *gin.Context) {
	submission := new(Submission)

	if err := ctx.BindJSON(submission); err != nil {

		ctx.IndentedJSON(http.StatusUnprocessableEntity, err)

	} else if c1, err := cr.repositoryAPI.CreateSubmission(submission); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusCreated, c1)
	}

}

// GetSubmission f
func (cr *SubmissionRouter) GetSubmission(ctx *gin.Context) {
	submissionID := ctx.Param("id")
	if submission, err := cr.repositoryAPI.GetSubmission(submissionID); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, submission)
	}
}

// GetSubmissions f
func (cr *SubmissionRouter) GetSubmissions(ctx *gin.Context) {
	if submissions, err := cr.repositoryAPI.GetSubmissions(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, submissions)
	}
}

