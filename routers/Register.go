package routers

import (
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/gin-gonic/gin"
)

// RegisterService s
type RegisterService struct {
	Router        *gin.Engine
	submissionAPI types.SubmissionAPI
}

// NewRegisterService f
func NewRegisterService(router *gin.Engine, submissionAPI types.SubmissionAPI) types.RegisterAPI {

	service := RegisterService{router, submissionAPI}
	return service

}

// Register f
func (rs RegisterService) Register() error {
	rs.Router.POST("/submissions", rs.submissionAPI.CreateSubmission)
	rs.Router.GET("/submissions", rs.submissionAPI.GetSubmissions)
	rs.Router.GET("/submissions/:id", rs.submissionAPI.GetSubmission)

	return nil
}
