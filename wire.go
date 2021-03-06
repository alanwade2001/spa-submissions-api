//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitialiseServerAPI() ServerAPI {
	wire.Build(
		gin.Default,
		NewMongoService,
		NewXMLParserAPI,
		NewGroupHeaderMapper,
		NewPaymentInformationMapper,
		NewInitiationMapper,
		NewCustomerService,
		NewValidator,
		NewInitiationService,
		NewSubmissionService,
		NewUserService,
		NewSubmissionRouter,
		NewRegisterService,
		NewConfigService,
		NewServer,
	)

	return &Server{}
}
