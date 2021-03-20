//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/alanwade2001/spa-submissions-api/repositories"
	"github.com/alanwade2001/spa-submissions-api/routers"
	"github.com/alanwade2001/spa-submissions-api/services"
	"github.com/alanwade2001/spa-submissions-api/types"
)

func InitialiseServerAPI() types.ServerAPI {
	wire.Build(
		gin.Default,
		services.NewMongoIdGenerator,
		repositories.NewMongoService,
		services.NewMessageService,
		services.NewXMLParserAPI,
		services.NewGroupHeaderMapper,
		services.NewPaymentInformationMapper,
		services.NewInitiationMapper,
		services.NewCustomerService,
		services.NewValidator,
		services.NewInitiationService,
		services.NewSubmissionService,
		services.NewUserService,
		routers.NewSubmissionRouter,
		routers.NewRegisterService,
		services.NewConfigService,
		NewServer,
	)

	return &Server{}
}

func InitialiseMockedServerAPI() types.ServerAPI {
	wire.Build(
		gin.Default,
		services.NewMongoIdGenerator,
		repositories.NewMongoService,
		services.NewMessageService,
		services.NewXMLParserAPI,
		services.NewGroupHeaderMapper,
		services.NewPaymentInformationMapper,
		services.NewInitiationMapper,
		services.NewMockCustomerService,
		services.NewValidator,
		services.NewInitiationService,
		services.NewSubmissionService,
		services.NewUserService,
		routers.NewSubmissionRouter,
		routers.NewRegisterService,
		services.NewConfigService,
		NewServer,
	)

	return &Server{}
}
