package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	hfiledata "service-worker-sqs-s3-postgres/entrypoints/controllers/filedata"
	hmetadata "service-worker-sqs-s3-postgres/entrypoints/controllers/metadata"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const rootPrefix = "/service-worker-sqs-s3-postgres"

// Server is an instance of Http Server for Rest endpoints.
type Server struct {
	server    *echo.Echo
	startedAt time.Time
	port      int
}

// NewServer creates an instance of Http Server.
func NewServer(port int, ec *hfiledata.FileDataController, mc *hmetadata.MetaDataController) *Server {
	e := echo.New()

	// middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server := &Server{server: e, port: port}

	// prefix
	path := e.Group(rootPrefix)

	// filedata
	path.GET("/s3/filedata/:id", ec.GetID)

	// metadata
	path.GET("/s3/metadata/:trackid", mc.GetID)

	return server
}

// Start runs a http server.
func (s *Server) Start() error {
	port := fmt.Sprintf(":%v", s.port)

	if err := s.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

// Stop stops an http server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Shutdown: %w", err)
	}

	return nil
}
