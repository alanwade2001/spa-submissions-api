// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitialiseServerAPI() ServerAPI {
	engine := gin.Default()
	repositoryAPI := NewMongoService()
	xmlParserAPI := NewXMLParserAPI()
	groupHeaderMapperAPI := NewGroupHeaderMapper()
	paymentInformationMapperAPI := NewPaymentInformationMapper()
	initiationMapperAPI := NewInitiationMapper(groupHeaderMapperAPI, paymentInformationMapperAPI)
	initiationAPI := NewInitiationService(xmlParserAPI, initiationMapperAPI)
	validatorAPI := NewValidator()
	submissionServiceAPI := NewSubmissionService(repositoryAPI, initiationAPI, validatorAPI)
	submissionAPI := NewSubmissionRouter(repositoryAPI, submissionServiceAPI)
	registerAPI := NewRegisterService(engine, submissionAPI)
	configAPI := NewConfigService()
	serverAPI := NewServer(engine, registerAPI, configAPI)
	return serverAPI
}
