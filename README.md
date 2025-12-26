# Kex (Knowledge Exchange)

**English** | [日本語](./README.ja.md)


**Kex** is a lightweight Document Librarian and MCP (Model Context Protocol) Server designed to manage and serve coding guidelines and project documentation. It helps AI agents (and humans) access the right documentation at the right time.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)

## Features

-   **MCP Server**: Exposes documents via the Model Context Protocol (JSON-RPC), allowing AI agents to "read" and "search" your documentation.
-   **Structure Enforcement**: Enforces a strict schema for documentation using Frontmatter (ID, Title, Status, Keywords).
-   **Lifecycle Management**: Distinguishes between `draft` and `adopted` documents.
    -   **Draft**: Work-in-progress. Errors are warnings.
    -   **Adopted**: Source of truth. Errors avoid startup.
-   **CLI Tooling**:
    -   `init`: Scaffolds a new Knowledge Base.
    -   `check`: Validates all documents and reports integrity issues.
    -   `start`: Starts the MCP server over Stdio.

## Installation

```bash
# Clone the repository
git clone https://github.com/mew-ton/kex.git
cd kex

# Build and Install
go install ./cmd/kex
```

## Getting Started

### 1. Initialize a Repository

Go to your project root and run:

```bash
kex init
```

This will create a `contents/` directory and a `.kex.yaml` configuration file.

### 2. Write Documentation

Create markdown files in the `contents/` directory. Each file **must** have valid YAML frontmatter:

```markdown
---
id: my-guideline
title: Coding Standards
description: Guidelines for Go development
status: adopted
keywords: [go, style, lint]
---

# Coding Standards

Write your content here...
```

**Status Types:**
-   `draft`: Ignored by the server by default. Useful for work-in-progress.
-   `adopted`: Active and indexed. Must pass all checks.

### 3. Validate Documents

Run the check command to verify your documents:

```bash
kex check
```

This will report:
-   Missing frontmatter
-   Filename vs ID mismatches
-   Parsing errors

### 4. Start the Server

Start the MCP server to allow AI connections:

```bash
kex start
```

*Note: This starts an interactive JSON-RPC session on Stdio. It is meant to be run by an MCP Client (like Claude Desktop).*

## MCP Tools Provided

Kex exposes the following tools to connected AI agents:

-   `search_documents(keywords: string[])`: Find documents matching specific keywords.
-   `read_document(id: string)`: Retrieve the full content of a document by its ID.

## Configuration

`.kex.yaml`:

```yaml
root: contents # Directory containing markdown files
```

## License

MIT