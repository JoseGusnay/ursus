package mcp

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// UrsusMCPServer handles MCP requests and maps them to Ursus services.
type UrsusMCPServer struct {
	server         *server.MCPServer
	service        *service.MemoryService
	sessionSvc     *service.SessionService
	suggestTopicUC *usecase.SuggestTopicUseCase
	timelineUC     *usecase.GetTimelineUseCase
	summarizeUC    *usecase.SummarizeSessionUseCase
	getDetailUC    *usecase.GetMemoryDetailUseCase
	passiveUC      *usecase.PassiveCaptureUseCase
	statsUC        *usecase.GetStatsUseCase
}

// NewUrsusMCPServer creates and configures a new MCP server for Ursus.
func NewUrsusMCPServer(svc *service.MemoryService, ss *service.SessionService, stuc *usecase.SuggestTopicUseCase, tuc *usecase.GetTimelineUseCase, sumuc *usecase.SummarizeSessionUseCase, gduc *usecase.GetMemoryDetailUseCase, puc *usecase.PassiveCaptureUseCase, statsuc *usecase.GetStatsUseCase) *UrsusMCPServer {
	s := server.NewMCPServer(
		"Ursus Memory System",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	ursusServer := &UrsusMCPServer{
		server:         s,
		service:        svc,
		sessionSvc:     ss,
		suggestTopicUC: stuc,
		timelineUC:     tuc,
		summarizeUC:    sumuc,
		getDetailUC:    gduc,
		passiveUC:      puc,
		statsUC:        statsuc,
	}

	ursusServer.registerTools()
	return ursusServer
}

func (s *UrsusMCPServer) registerTools() {
	// Add Memory Tool
	addTool := mcp.NewTool("add_memory",
		mcp.WithDescription("Add a new memory to Ursus"),
		mcp.WithString("content", mcp.Description("Content of the memory to store"), mcp.Required()),
		mcp.WithString("metadata", mcp.Description("Optional metadata associated with the memory")),
		mcp.WithString("topic_key", mcp.Description("Optional key for topic-based upserts")),
		mcp.WithString("scope", mcp.Description("Scope of the memory: 'project' (default) or 'personal'")),
	)

	// Search Memory Tool
	searchTool := mcp.NewTool("search_memory",
		mcp.WithDescription("Search for memories in Ursus"),
		mcp.WithString("query", mcp.Description("Search query"), mcp.Required()),
	)

	// Session Start Tool
	sessionStartTool := mcp.NewTool("session_start",
		mcp.WithDescription("Start a new work session"),
		mcp.WithString("title", mcp.Description("Title of the session"), mcp.Required()),
	)

	// Session End Tool
	sessionEndTool := mcp.NewTool("session_end",
		mcp.WithDescription("End the active work session"),
	)

	// Suggest Topic Tool
	suggestTopicTool := mcp.NewTool("suggest_topic",
		mcp.WithDescription("Suggest relevant topics based on stored memories"),
	)

	// Timeline Tool
	timelineTool := mcp.NewTool("get_timeline",
		mcp.WithDescription("Get a chronological timeline of stored memories"),
	)

	// Summarize Session Tool
	summarizeTool := mcp.NewTool("summarize_session",
		mcp.WithDescription("Summarize a specific work session or the current one"),
		mcp.WithString("session_id", mcp.Description("Optional session ID to summarize")),
	)

	// Detail Tool (3-Layer Pattern: Layer 3)
	detailTool := mcp.NewTool("get_memory_detail",
		mcp.WithDescription("Get full untruncated content and metadata of a specific memory"),
		mcp.WithString("id", mcp.Description("ID of the memory to retrieve")),
	)

	// Update Tool
	updateTool := mcp.NewTool("update_memory",
		mcp.WithDescription("Update an existing memory with new content or metadata"),
		mcp.WithString("id", mcp.Description("ID of the memory to update")),
		mcp.WithString("content", mcp.Description("New content")),
		mcp.WithString("metadata", mcp.Description("New metadata")),
	)

	// Delete Tool
	deleteTool := mcp.NewTool("delete_memory",
		mcp.WithDescription("Delete a memory by its ID"),
		mcp.WithString("id", mcp.Description("ID of the memory to delete")),
	)

	s.server.AddTool(addTool, s.handleAddMemory)
	s.server.AddTool(searchTool, s.handleSearchMemory)
	s.server.AddTool(sessionStartTool, s.handleSessionStart)
	s.server.AddTool(sessionEndTool, s.handleSessionEnd)
	s.server.AddTool(suggestTopicTool, s.handleSuggestTopic)
	s.server.AddTool(timelineTool, s.handleGetTimeline)
	s.server.AddTool(summarizeTool, s.handleSummarizeSession)
	s.server.AddTool(detailTool, s.handleGetMemoryDetail)
	s.server.AddTool(updateTool, s.handleUpdateMemory)
	s.server.AddTool(deleteTool, s.handleDeleteMemory)

	// Passive Capture Tool
	passiveTool := mcp.NewTool("passive_capture",
		mcp.WithDescription("Extract and save learnings automatically from text using markers like '### Aprendizajes' or '<learning>' tags"),
		mcp.WithString("text", mcp.Description("The full text to scan for learnings"), mcp.Required()),
	)
	s.server.AddTool(passiveTool, s.handlePassiveCapture)

	// Stats Tool
	statsTool := mcp.NewTool("mem_stats",
		mcp.WithDescription("Get a comprehensive report of Ursus memory statistics"),
	)
	s.server.AddTool(statsTool, s.handleGetStats)
}

func (s *UrsusMCPServer) handleAddMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content := request.GetString("content", "")
	metadata := request.GetString("metadata", "")

	topicKey := request.GetString("topic_key", "")
	scope := request.GetString("scope", entity.ScopeProject)

	mem, err := s.service.Store(ctx, content, metadata, topicKey, scope, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Store error: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Memory saved successfully! ID: %s", mem.ID)), nil
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

func (s *UrsusMCPServer) handleSessionStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title := request.GetString("title", "")
	if title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}
	sess, err := s.sessionSvc.Start(ctx, title)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Session start error: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Session started: %s (ID: %s)", sess.Title, sess.ID)), nil
}

func (s *UrsusMCPServer) handleSuggestTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topics, err := s.suggestTopicUC.Execute(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Suggest topic error: %v", err)), nil
	}

	if len(topics) == 0 {
		return mcp.NewToolResultText("No enough context to suggest topics yet."), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Suggested topics based on your memories: %s", strings.Join(topics, ", "))), nil
}

func (s *UrsusMCPServer) handleSessionEnd(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if err := s.sessionSvc.End(ctx); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Session end error: %v", err)), nil
	}
	return mcp.NewToolResultText("Session ended successfully"), nil
}

func (s *UrsusMCPServer) handleGetTimeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	days, err := s.timelineUC.Execute(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Timeline error: %v", err)), nil
	}

	if len(days) == 0 {
		return mcp.NewToolResultText("No memories found."), nil
	}

	var response strings.Builder
	response.WriteString("Ursus Memory Timeline:\n")
	for _, day := range days {
		response.WriteString(fmt.Sprintf("\n--- %s ---\n", day.Date.Format("2006-01-02")))
		for _, m := range day.Memories {
			timeStr := m.CreatedAt.Format("15:04")
			response.WriteString(fmt.Sprintf("  [%s] %s\n", timeStr, m.Content))
		}
	}

	return mcp.NewToolResultText(response.String()), nil
}

func (s *UrsusMCPServer) handleSummarizeSession(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sessionID := request.GetString("session_id", "")
	review, err := s.summarizeUC.Execute(ctx, sessionID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Summarize error: %v", err)), nil
	}

	return mcp.NewToolResultText(review.Summary), nil
}

func (s *UrsusMCPServer) handleGetMemoryDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	if id == "" {
		return mcp.NewToolResultError("ID is required"), nil
	}

	mem, err := s.getDetailUC.Execute(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Get detail error: %v", err)), nil
	}
	if mem == nil {
		return mcp.NewToolResultError("Memory not found"), nil
	}

	response := fmt.Sprintf("ID: %s\nContent: %s\nMetadata: %s\nTopicKey: %s\nRevision: %d\nCreated: %s",
		mem.ID, mem.Content, mem.Metadata, mem.TopicKey, mem.RevisionCount, mem.CreatedAt.Format("2006-01-02 15:04"))

	return mcp.NewToolResultText(response), nil
}

func (s *UrsusMCPServer) handleUpdateMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	content := request.GetString("content", "")
	metadata := request.GetString("metadata", "")

	if id == "" {
		return mcp.NewToolResultError("ID is required"), nil
	}

	mem, err := s.service.GetByID(ctx, id)
	if err != nil || mem == nil {
		return mcp.NewToolResultError("Memory not found"), nil
	}

	if content != "" {
		mem.Content = content
	}
	if metadata != "" {
		mem.Metadata = metadata
	}

	if err := s.service.Update(ctx, mem); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Update error: %v", err)), nil
	}

	return mcp.NewToolResultText("Memory updated successfully"), nil
}

func (s *UrsusMCPServer) handlePassiveCapture(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text := request.GetString("text", "")
	if text == "" {
		return mcp.NewToolResultError("text is required"), nil
	}

	mems, err := s.passiveUC.Execute(ctx, text)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Passive capture error: %v", err)), nil
	}

	if len(mems) == 0 {
		return mcp.NewToolResultText("No learnings found to capture."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Successfully captured %d learnings:\n", len(mems)))
	for _, m := range mems {
		sb.WriteString(fmt.Sprintf("- %s\n", m.Content))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (s *UrsusMCPServer) handleGetStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	report, err := s.statsUC.Execute(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Stats error: %v", err)), nil
	}

	var sb strings.Builder
	sb.WriteString("Ursus Memory Statistics:\n")
	sb.WriteString(fmt.Sprintf("- Total Memories: %d\n", report.TotalMemories))
	sb.WriteString(fmt.Sprintf("- Total Sessions: %d\n", report.TotalSessions))
	sb.WriteString(fmt.Sprintf("- Total Prompts Logged: %d\n", report.TotalPrompts))
	
	if len(report.TopTopics) > 0 {
		sb.WriteString("\nTop Topics:\n")
		for _, t := range report.TopTopics {
			sb.WriteString(fmt.Sprintf("  - %s\n", t))
		}
	}

	if len(report.Last7DaysActivity) > 0 {
		sb.WriteString("\nActivity (Last 7 Days):\n")
		// Sort keys for consistent output
		var dates []string
		for d := range report.Last7DaysActivity {
			dates = append(dates, d)
		}
		sort.Strings(dates)
		for _, d := range dates {
			sb.WriteString(fmt.Sprintf("  - %s: %d memories\n", d, report.Last7DaysActivity[d]))
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (s *UrsusMCPServer) handleDeleteMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	if id == "" {
		return mcp.NewToolResultError("ID is required"), nil
	}

	if err := s.service.Delete(ctx, id); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Delete error: %v", err)), nil
	}

	return mcp.NewToolResultText("Memory deleted successfully"), nil
}

// Serve starts the MCP server on stdio.
func (s *UrsusMCPServer) Serve() error {
	return server.ServeStdio(s.server)
}
