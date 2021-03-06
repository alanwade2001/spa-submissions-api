package routers

import (
	"net/http"

	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/gin-gonic/gin"
)

// SubmissionRouter s
type SubmissionRouter struct {
	userAPI    types.UserAPI
	serviceAPI types.SubmissionServiceAPI
}

// NewSubmissionRouter f
func NewSubmissionRouter(userAPI types.UserAPI, serviceAPI types.SubmissionServiceAPI) types.SubmissionAPI {

	submissionAPI := SubmissionRouter{userAPI, serviceAPI}

	return &submissionAPI
}

// CreateSubmission f
func (cr *SubmissionRouter) CreateSubmission(ctx *gin.Context) {

	if user, err := cr.userAPI.Find(ctx); err != nil {
		ctx.String(http.StatusUnauthorized, err.Error())
	} else if submission, err := cr.serviceAPI.CreateSubmission(ctx.Request.Body, *user); err != nil {
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
	} else if submission == nil {
		ctx.Status(http.StatusNotFound)
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
