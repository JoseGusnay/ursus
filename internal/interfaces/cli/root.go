package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/infrastructure/api"
	"github.com/JoseGusnay/ursus/internal/interfaces/tui"
	"github.com/spf13/cobra"
)

var (
	svc         *service.MemoryService
	sessionSvc  *service.SessionService
	gitSvc      *service.GitService
	syncUC      *usecase.SyncMemoriesUseCase
	suggestUC   *usecase.SuggestTopicUseCase
	timelineUC  *usecase.GetTimelineUseCase
	summarizeUC *usecase.SummarizeSessionUseCase
	getDetailUC *usecase.GetMemoryDetailUseCase
	statsUC     *usecase.GetStatsUseCase
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ursus",
	Short: "Ursus is a tool for managing persistent context and memory",
	Long:  `Ursus is a CLI and MCP server that allows AI agents to remember context across sessions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(s *service.MemoryService, ss *service.SessionService, suc *usecase.SyncMemoriesUseCase, stuc *usecase.SuggestTopicUseCase, tuc *usecase.GetTimelineUseCase, sumuc *usecase.SummarizeSessionUseCase, gduc *usecase.GetMemoryDetailUseCase, stauc *usecase.GetStatsUseCase) {
	svc = s
	sessionSvc = ss
	syncUC = suc
	suggestUC = stuc
	timelineUC = tuc
	summarizeUC = sumuc
	getDetailUC = gduc
	statsUC = stauc
	gitSvc = service.NewGitService()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the Ursus REST API server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		addr := ":" + port
		server := api.NewRESTServer(svc)
		fmt.Printf("Ursus REST API starting on http://localhost%s\n", addr)
		if err := server.ListenAndServe(addr); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync memories with current directory (Git-friendly Gzip Chunks)",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		
		fmt.Println("Importing chunks from .ursus/...")
		importUC := usecase.NewImportChunksUseCase(svc.Repository())
		if err := importUC.Execute(context.Background(), cwd); err != nil {
			fmt.Printf("Import error: %v\n", err)
		}

		fmt.Println("Exporting local memories to compressed chunks...")
		exportUC := usecase.NewExportChunksUseCase(svc.Repository())
		user := os.Getenv("USER")
		if user == "" {
			user = "anonymous"
		}
		if err := exportUC.Execute(context.Background(), cwd, user); err != nil {
			fmt.Printf("Export error: %v\n", err)
			return
		}

		fmt.Println("Sync completed successfully!")
	},
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open the Ursus TUI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Start(svc); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add [content]",
	Short: "Add a new memory",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		metadata, _ := cmd.Flags().GetString("metadata")
		topic, _ := cmd.Flags().GetString("topic")
		content := args[0]
		mem, err := svc.Store(context.Background(), content, metadata, topic, "", "")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Memory saved successfully! ID: %s\n", mem.ID)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search memories",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		results, err := svc.Search(context.Background(), args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(results) == 0 {
			fmt.Println("No memories found.")
			return
		}
		for _, r := range results {
			fmt.Printf("[%s] ID: %s\n%s\n---\n", r.CreatedAt.Format("2006-01-02"), r.ID, r.Content)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all memories",
	Run: func(cmd *cobra.Command, args []string) {
		results, err := svc.List(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		for _, r := range results {
			fmt.Printf("[%s] ID: %s\n%s\n---\n", r.CreatedAt.Format("2006-01-02"), r.ID, r.Content)
		}
	},
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage work sessions",
}

var sessionStartCmd = &cobra.Command{
	Use:   "start [title]",
	Short: "Start a new work session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := sessionSvc.Start(context.Background(), args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Session started: %s (ID: %s)\n", s.Title, s.ID)
	},
}

var sessionEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End the active work session",
	Run: func(cmd *cobra.Command, args []string) {
		if err := sessionSvc.End(context.Background()); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Session ended.")
	},
}

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggest relevant topics based on stored memories",
	Run: func(cmd *cobra.Command, args []string) {
		topics, err := suggestUC.Execute(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(topics) == 0 {
			fmt.Println("No enough context to suggest topics yet.")
			return
		}
		fmt.Printf("Suggested topics: %s\n", strings.Join(topics, ", "))
	},
}

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show a chronological timeline of memories",
	Run: func(cmd *cobra.Command, args []string) {
		days, err := timelineUC.Execute(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(days) == 0 {
			fmt.Println("No memories found.")
			return
		}

		for _, day := range days {
			fmt.Printf("\n--- %s ---\n", day.Date.Format("2006-01-02"))
			for _, m := range day.Memories {
				timeStr := m.CreatedAt.Format("15:04")
				fmt.Printf("  [%s] %s\n", timeStr, m.Content)
			}
		}
	},
}

var reviewCmd = &cobra.Command{
	Use:   "review [session-id]",
	Short: "Review and summarize a work session",
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := ""
		if len(args) > 0 {
			sessionID = args[0]
		}

		review, err := summarizeUC.Execute(context.Background(), sessionID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("\n=== SESSION REVIEW: %s ===\n", review.Session.Title)
		fmt.Printf("Started: %s\n", review.Session.StartTime.Format("2006-01-02 15:04"))
		fmt.Println("-------------------------------------------")
		fmt.Println(review.Summary)
		fmt.Println("-------------------------------------------")
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show memory statistics",
	Run: func(cmd *cobra.Command, args []string) {
		report, err := statsUC.Execute(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("\n=== URSUS MEMORY STATS ===")
		fmt.Printf("Total Memories: %d\n", report.TotalMemories)
		fmt.Printf("Total Sessions: %d\n", report.TotalSessions)
		fmt.Printf("Total Prompts:  %d\n", report.TotalPrompts)
		fmt.Println("\nTop Topics:")
		for _, t := range report.TopTopics {
			fmt.Printf(" - %s\n", t)
		}
		fmt.Println("\nActivity (Last 7 Days):")
		for date, count := range report.Last7DaysActivity {
			fmt.Printf(" %s: %d items\n", date, count)
		}
	},
}

var detailCmd = &cobra.Command{
	Use:   "detail [id]",
	Short: "Show full details of a memory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mem, err := getDetailUC.Execute(context.Background(), args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if mem == nil {
			fmt.Println("Memory not found.")
			return
		}
		fmt.Printf("\n=== MEMORY DETAIL ===\n")
		fmt.Printf("ID:       %s\n", mem.ID)
		fmt.Printf("Scope:    %s\n", mem.Scope)
		fmt.Printf("Topic:    %s\n", mem.TopicKey)
		fmt.Printf("Duplicates: %d\n", mem.DuplicateCount)
		fmt.Printf("Revisions:  %d\n", mem.RevisionCount)
		fmt.Printf("Created:  %s\n", mem.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("Last Seen: %s\n", mem.LastSeenAt.Format("2006-01-02 15:04"))
		fmt.Printf("\nCONTENT:\n%s\n", mem.Content)
		if mem.Metadata != "" {
			fmt.Printf("\nMETADATA: %s\n", mem.Metadata)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(searchCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(tuiCmd)
	RootCmd.AddCommand(syncCmd)
	RootCmd.AddCommand(apiCmd)
	RootCmd.AddCommand(sessionCmd)
	RootCmd.AddCommand(suggestCmd)
	RootCmd.AddCommand(timelineCmd)
	RootCmd.AddCommand(reviewCmd)
	RootCmd.AddCommand(statsCmd)
	RootCmd.AddCommand(detailCmd)

	sessionCmd.AddCommand(sessionStartCmd)
	sessionCmd.AddCommand(sessionEndCmd)

	apiCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	addCmd.Flags().StringP("metadata", "m", "", "Optional metadata for the memory")
	addCmd.Flags().StringP("topic", "t", "", "Optional topic key for upserts")
}
