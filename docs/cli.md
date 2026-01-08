# CLI Reference

## `kex init`

Initializes a new Kex repository.

```bash
kex init [options]
```

- **Interactive Mode**: Prompts for Agent Type and Scopes (Coding, Documentation) if no flags provided.
- **Flags**: `--agent-type=<type>` (e.g., `general`, `claude`) to bypass prompts.
- Creates `.kex.yaml` with selected configuration.
- Creates `contents/` directory.
- Generates `AGENTS.md` (or `CLAUDE.md`) with dynamic content based on selected scopes.

## `kex check`

Validates the integrity of your documentation repository.

```bash
kex check [options] [path...]
```

You can specify multiple documentation sources (directories or URLs) to validate them as a unified set.

**Examples**:
```bash
kex check                           # Use .kex.yaml sources
kex check contents/coding           # Check specific dir
kex check contents/coding docs      # Check multiple dirs
```

checks for:
- Invalid YAML Frontmatter
- ID vs Filename mismatches
- Missing required fields
- Duplicate IDs (across all sources)

## `kex start`

Starts the MCP Server.

```bash
kex start [options] [path...]
```

**Argument Logic**:
Kex inspects each provided path to determine if it is a **Project Root** (contains `.kex.yaml` or `kex.json`) or a **Content Source**.

For detailed resolution rules and how to work with multiple sources, see [Composition](composition.md#resolution-logic).

**Examples**:
- **From Config**: `kex start .` (Loads `.kex.yaml` from current dir)
- **Ad-hoc**: `kex start ./contents` (Starts server with just ./contents)
- **Mixed**: `kex start . ./extra-docs`

**Flags**:
- `--log-file=<path>`: Write logs to a file instead of Stderr.

## `kex generate`

Generates a static site structure for remote hosting (GitHub Pages).

```bash
kex generate [options] [path...]
```

- Combines all specified sources into a single static site.
- Creates a `dist/` directory in the current working directory.
- Copies all "Adopted" documents and generates `kex.json`.

**Example**:
```bash
kex generate contents/coding contents/docs
```

## `kex update`

Updates the Kex system documentation and agent configuration in an existing repository.

```bash
kex update [options]
```

- **System Docs (`contents/documentation/kex/*`)**: Updates to match the current binary version (Overwrite).
- **`AGENTS.md`**: Safely updates guidelines between `<!-- kex: auto-update start -->` markers.
- See `.kex.yaml` configuration to customize behavior.
