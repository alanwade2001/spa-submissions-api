package main

import "github.com/gin-gonic/gin"

// NewServer f
func NewServer(router *gin.Engine, registerAPI RegisterAPI, configAPI ConfigAPI) ServerAPI {

	return &Server{router, registerAPI, configAPI}
}

// Server s
type Server struct {
	Router      *gin.Engine
	registerAPI RegisterAPI
	configAPI   ConfigAPI
}

// Run f
func (s *Server) Run() error {
	if err := s.configAPI.Load(); err != nil {
		return err
	}

	if err := s.registerAPI.Register(); err != nil {
		return err
	}

	return s.Router.Run()
}
