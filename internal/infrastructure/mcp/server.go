package mcp

import (
	"context"
	"fmt"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// UrsusMCPServer handles MCP requests and maps them to Ursus services.
type UrsusMCPServer struct {
	service *service.MemoryService
	mcp     *server.MCPServer
}

// NewUrsusMCPServer creates and configures a new MCP server for Ursus.
func NewUrsusMCPServer(svc *service.MemoryService) *UrsusMCPServer {
	s := server.NewMCPServer(
		"Ursus Memory System",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)

	ursusServer := &UrsusMCPServer{
		service: svc,
		mcp:     s,
	}

	ursusServer.registerTools()
	return ursusServer
}

func (s *UrsusMCPServer) registerTools() {
	// Tool: add_memory
	addTool := mcp.NewTool("add_memory",
		mcp.WithDescription("Saves a new memory or observation to Ursus"),
		mcp.WithString("content", mcp.Description("The main content of the memory"), mcp.Required()),
		mcp.WithString("metadata", mcp.Description("Optional context or tags")),
	)

	s.mcp.AddTool(addTool, s.handleAddMemory)

	// Tool: search_memory
	searchTool := mcp.NewTool("search_memory",
		mcp.WithDescription("Searches for existing memories in Ursus using keywords"),
		mcp.WithString("query", mcp.Description("The search term"), mcp.Required()),
	)

	s.mcp.AddTool(searchTool, s.handleSearchMemory)
}

func (s *UrsusMCPServer) handleAddMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content := request.GetString("content", "")
	if content == "" {
		return mcp.NewToolResultError("content must be provided and be a string"), nil
	}
	metadata := request.GetString("metadata", "")

	u, err := s.service.Store(ctx, content, metadata)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to save memory: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Memory saved successfully with ID: %s", u.ID)), nil
}

func (s *UrsusMCPServer) handleSearchMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")
	if query == "" {
		return mcp.NewToolResultError("query must be provided and be a string"), nil
	}

	results, err := s.service.Search(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No memories found for that query."), nil
	}

	response := "Found results:\n"
	for _, r := range results {
		response += fmt.Sprintf("- [%s] %s\n", r.CreatedAt.Format("2006-01-02"), r.Content)
	}

	return mcp.NewToolResultText(response), nil
}

// Serve starts the MCP server on stdio.
func (s *UrsusMCPServer) Serve() error {
	return server.ServeStdio(s.mcp)
}
