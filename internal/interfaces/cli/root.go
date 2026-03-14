package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/infrastructure/api"
	"github.com/JoseGusnay/ursus/internal/interfaces/tui"
	"github.com/spf13/cobra"
)

var (
	svc    *service.MemoryService
	gitSvc *service.GitService
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ursus",
	Short: "Ursus is a persistent memory system for AI agents",
	Long:  `Ursus is a CLI and MCP server that allows AI agents to remember context across sessions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(s *service.MemoryService) {
	svc = s
	gitSvc = service.NewGitService()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(searchCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(tuiCmd)
	RootCmd.AddCommand(syncCmd)
	RootCmd.AddCommand(apiCmd)
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

func init() {
	apiCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize memories with a Git repository",
}

var syncInitCmd = &cobra.Command{
	Use:   "init [remote-url]",
	Short: "Initialize Git sync with a remote repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := gitSvc.Init(args[0])
		if err != nil {
			fmt.Printf("Error: %v\nOutput: %s\n", err, out)
			return
		}
		fmt.Println("Git sync initialized successfully.")
	},
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push memories to the remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := gitSvc.Push()
		if err != nil {
			fmt.Printf("Error: %v\nOutput: %s\n", err, out)
			return
		}
		fmt.Println("Memories pushed to remote.")
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull memories from the remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := gitSvc.Pull()
		if err != nil {
			fmt.Printf("Error: %v\nOutput: %s\n", err, out)
			return
		}
		fmt.Println("Memories pulled from remote.")
	},
}

func init() {
	syncCmd.AddCommand(syncInitCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
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
		u, err := svc.Store(context.Background(), args[0], metadata)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Memory saved successfully! ID: %s\n", u.ID)
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

func init() {
	addCmd.Flags().StringP("metadata", "m", "", "Optional metadata for the memory")
}
