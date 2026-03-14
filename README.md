# Ursus 🐻

**Ursus** is a robust, persistent memory system designed for AI agents. Built with Go, SQLite, and the Model Context Protocol (MCP), it enables seamless context retention and long-term memory management for advanced AI coding agents and assistants.

## 🚀 Features

- **🧠 Long-term Memory**: Persistent storage for observations and context.
- **🔍 Semantic Search**: Full-text search (FTS5) for instant retrieval of relevant memories.
- **🤖 MCP Support**: Built-in Model Context Protocol server for integration with agents like Claude, Cursor, and Antigravity.
- **🖥️ Multi-Interface**:
  - **TUI**: Beautiful terminal-based visual interface.
  - **CLI**: Powerful command-line tools for manual management.
  - **REST API**: Standard HTTP interface for external integrations.
- **📄 Sync Support**: Portable memory via Git-based synchronization.
- **🏗️ Solid Architecture**: Built on Clean Architecture and SOLID principles.

## 🛠️ Quick Start

### Build from source
```bash
go build -o ursus ./cmd/ursus
```

### Usage
- **Interactive TUI**: `./ursus tui`
- **MCP Server**: `./ursus server`
- **REST API**: `./ursus api -p 8080`
- **Add Memory**: `./ursus add "Important context about this project"`
- **List Memories**: `./ursus list`

## 🏗️ Architecture

Ursus follows **Clean Architecture** patterns:
- **Domain**: Core entities and repository interfaces.
- **Application**: Business logic and use cases.
- **Infrastructure**: Implementation details (SQLite, MCP, REST, CLI).
- **Interfaces**: Entry points for user/agent interaction.

---
Developed by Jose Gusnay & Antigravity
