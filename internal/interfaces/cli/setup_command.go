package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Automate Ursus installation and agent integration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Ursus Setup Wizard")
		fmt.Println("-------------------")
		fmt.Println("1. Configure system PATH (Windows)")
		fmt.Println("2. Configure Claude Desktop")
		fmt.Println("3. Configure Cursor/Windsurf")
		fmt.Println("-------------------")
		fmt.Println("Run 'ursus setup [command]' for specific configuration.")
	},
}

var setupPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Add Ursus to system PATH (Windows)",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "windows" {
			fmt.Println("This command is only supported on Windows.")
			return
		}

		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			return
		}

		binDir := filepath.Dir(exePath)
		fmt.Printf("Adding %s to User PATH...\n", binDir)

		// Get current PATH
		out, err := exec.Command("powershell", "-Command", "[Environment]::GetEnvironmentVariable('Path', 'User')").Output()
		if err != nil {
			fmt.Printf("Error getting current PATH: %v\n", err)
			return
		}

		currentPath := string(out)
		if strings.Contains(currentPath, binDir) {
			fmt.Println("✅ Directory already in PATH.")
			return
		}

		newPath := binDir + ";" + strings.TrimSpace(currentPath)
		setx := exec.Command("setx", "Path", newPath)
		if err := setx.Run(); err != nil {
			fmt.Printf("Error setting PATH: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully added to PATH! Please restart your terminal.")
	},
}

var setupClaudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Add Ursus to Claude Desktop configuration",
	Run: func(cmd *cobra.Command, args []string) {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			fmt.Println("Could not find APPDATA directory.")
			return
		}

		configPath := filepath.Join(appData, "Claude", "claude_desktop_config.json")
		fmt.Printf("Looking for config at: %s\n", configPath)

		// Create directory if it doesn't exist
		_ = os.MkdirAll(filepath.Dir(configPath), 0755)

		var config map[string]interface{}
		data, err := os.ReadFile(configPath)
		if err == nil {
			err = json.Unmarshal(data, &config)
			if err != nil {
				fmt.Println("Error parsing existing config, creating new one.")
				config = make(map[string]interface{})
			}
		} else {
			config = make(map[string]interface{})
		}

		if _, ok := config["mcpServers"]; !ok {
			config["mcpServers"] = make(map[string]interface{})
		}

		mcpServers := config["mcpServers"].(map[string]interface{})
		exePath, _ := os.Executable()

		mcpServers["ursus"] = map[string]interface{}{
			"command": exePath,
			"args":    []string{"mcp"},
		}

		newData, _ := json.MarshalIndent(config, "", "  ")
		err = os.WriteFile(configPath, newData, 0644)
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully added Ursus to Claude Desktop!")
	},
}

func init() {
	setupCmd.AddCommand(setupPathCmd)
	setupCmd.AddCommand(setupClaudeCmd)
	// Cursor is usually manual but we can provide instructions
	setupCmd.AddCommand(&cobra.Command{
		Use:   "cursor",
		Short: "Get instructions for Cursor integration",
		Run: func(cmd *cobra.Command, args []string) {
			exePath, _ := os.Executable()
			fmt.Println("\n🤖 Cursor / Windsurf Integration:")
			fmt.Println("1. Open Cursor Settings -> Features -> MCP.")
			fmt.Println("2. Click '+ Add New MCP Server'.")
			fmt.Printf("3. Command: %s mcp\n", exePath)
			fmt.Println("4. Name: Ursus")
			fmt.Println("5. Type: command")
		},
	})
}
