# Core Concepts

## Philosophy

Kex is designed to be a **Librarian**, not a Judge.

- **Dynamic yet Efficient Context**: AI should fetch context dynamically when needed, but the effort (and context window) to do so should be minimal.
- **Clean System Prompts**: Coding guidelines should **not** be hardcoded in project prompts (e.g., `.cursorrules`). They belong in the repository (or a central knowledge base), indexed by Kex.
- **Structure is strict, meaning is human-owned**: We enforce schema (Frontmatter), but content is up to you.
- **Centralized Knowledge**: The goal is to separate domain-agnostic knowledge from project specific prompts. Kex allows you to centralize these guidelines and have the Librarian serve them, preventing the need to "scold" AI for missing rules.

## Document Lifecycle

Kex distinguishes between two states for documents:

- **Draft**: Work-in-progress. Ignored by the server by default (unless configured otherwise). Errors are treated as warnings.
- **Adopted**: The source of truth. Active and indexed. Must pass all validation checks. Errors prevent the server from starting.

## Architecture

Kex operates as a local **MCP (Model Context Protocol) Server**.
It runs as a single binary that indexes markdown files located in a specified directory (local or remote).

It serves these documents over JSON-RPC to supported AI clients (Claude, VSCode, etc).
