# Kex

**English** | [日本語](./README.ja.md)

**Kex** is a **keyword-indexed knowledge exchange** implemented as a **document librarian tool (MCP / Skills Management)**. Designed to aid "Vibe Coding", it operates as a local **MCP Server** that helps AI agents (and humans) access the right documentation at the right time from your repository.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-blue)

## Quick Start

See **[Installation](docs/installation.md)** for binary downloads and `go install` instructions.

### 1. Initialize
```bash
kex init
```

### 2. Start (Local)
```bash
kex start .
```

### 3. Start (Remote / GitHub Pages)
```bash
kex start https://my-org.github.io/guidelines/
```

## Documentation

Full documentation is available in the `docs/` directory:

- **[Features](docs/feature-index.md)**: Keywords and Scopes explanations.
- **[Core Concepts](docs/concept.md)**: Philosophy, Architecture, and Draft/Adopted status.
- **[CLI Reference](docs/cli.md)**: Usage of `init`, `check`, `start`, and `generate`.
- **[Configuration](docs/configuration.md)**: `.kex.yaml` reference.
- **[MCP Tools](docs/feature-mcp.md)**: Tools exposed to AI agents and client configuration.
- **[AI Skills](docs/feature-skills.md)**: Dynamic knowledge generation for AI Agents.
- **[Writing Documentation](docs/documentation.md)**: Frontmatter schema and content guidelines.
- **[Best Practices](docs/best-practice.md)**: How to structure your knowledge base.

## Use Cases

- **[Monorepo / Local](docs/usecase-monorepo.md)**: For documentation that lives with the code.
- **[Central Repository](docs/usecase-central-repo.md)**: For organization-wide guidelines shared via GitHub Pages.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT