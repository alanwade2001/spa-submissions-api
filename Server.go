package main

import (
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/gin-gonic/gin"
)

// NewServer f
func NewServer(router *gin.Engine, registerAPI types.RegisterAPI, configAPI types.ConfigAPI) types.ServerAPI {

	return &Server{router, registerAPI, configAPI}
}

// Server s
type Server struct {
	Router      *gin.Engine
	registerAPI types.RegisterAPI
	configAPI   types.ConfigAPI
}

// Run f
func (s *Server) Run() error {
	if err := s.configAPI.Load("."); err != nil {
		return err
	}

	if err := s.registerAPI.Register(); err != nil {
		return err
	}

	return s.Router.Run()
}
