// Package mcp provides MCP (Model Context Protocol) server for prism-roadmap.
package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/prism-roadmap/rmi"
)

// Server provides MCP tools for roadmap management.
type Server struct {
	services map[string]*rmi.Service
}

// NewServer creates a new MCP server.
func NewServer() *Server {
	return &Server{
		services: make(map[string]*rmi.Service),
	}
}

// RegisterTools registers all roadmap tools with the MCP server.
func (s *Server) RegisterTools(server *mcp.Server) {
	s.RegisterRMITools(server)
}

// getService returns a cached service for the given file path, creating one if needed.
func (s *Server) getService(filePath string) (*rmi.Service, error) {
	if svc, ok := s.services[filePath]; ok {
		return svc, nil
	}

	svc, err := rmi.NewServiceFromFile(filePath)
	if err != nil {
		return nil, err
	}
	s.services[filePath] = svc
	return svc, nil
}

// getOrCreateService returns a service, creating a new one if the file doesn't exist.
func (s *Server) getOrCreateService(filePath string) *rmi.Service {
	if svc, ok := s.services[filePath]; ok {
		return svc
	}

	svc, err := rmi.NewServiceFromFile(filePath)
	if err != nil {
		// Create new service
		svc = rmi.NewService()
	}
	s.services[filePath] = svc
	return svc
}

// Run starts the MCP server.
func (s *Server) Run(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "prism-roadmap",
		Version: "0.1.0",
	}, nil)

	s.RegisterTools(server)

	return server.Run(ctx, &mcp.StdioTransport{})
}
