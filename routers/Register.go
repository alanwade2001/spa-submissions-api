package routers

import (
	"bytes"
	"io/ioutil"

	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
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

	if err := rs.RegisterLogRequest(); err != nil {
		return nil
	}

	if err := rs.RegisterRoutes(); err != nil {
		return nil
	}

	return nil
}

// RegisterLogRequest f
func (rs RegisterService) RegisterLogRequest() error {
	rs.Router.Use(func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		klog.Infoln(string(body))

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	})

	return nil
}

// RegisterRoutes f
func (rs RegisterService) RegisterRoutes() error {
	rs.Router.POST("/submissions", rs.submissionAPI.CreateSubmission)
	rs.Router.GET("/submissions", rs.submissionAPI.GetSubmissions)
	rs.Router.GET("/submissions/:id", rs.submissionAPI.GetSubmission)

	return nil
}
