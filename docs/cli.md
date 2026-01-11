# CLI Reference

## `kex init`

Initializes a new Kex repository.

```bash
kex init [options]
```

- **Interactive Mode**: Prompts for AI Agents (Antigravity, Cursor, Claude) and Scopes (Coding, Documentation) if no flags provided.
- **Flags**:
    - `--agents=<agent>`: Select AI agents (can be repeated).
    - `--scopes=<scope>`: Select scopes (can be repeated).
- Creates `.kex.yaml` with selected configuration.
- Creates `contents/` directory.
- Generates Agent Rule Files (e.g. `.antigravity/rules/kex-coding.md`) based on selections.

## `kex check`

Validates the integrity of your documentation repository.

```bash
kex check
```

Checks for:
- Invalid YAML Frontmatter
- ID vs Filename mismatches
- Missing required fields
- Duplicate IDs

## `kex start`

Starts the MCP Server.

```bash
kex start [options] [path|url]
```

- **Local Mode**: `kex start ./my-docs`
- **Remote Mode**: `kex start https://example.com/docs/` (Must contain `kex.json`)
- **Flags**:
    - `--log-file=<path>`: Write logs to a file instead of Stderr.
    - `--root=<path>`: Specify the guidelines content root (Local mode).



## `kex generate`

Generates a static site structure for remote hosting (GitHub Pages).

```bash
kex generate [output-dir]
```

- Validates all "Adopted" documents.
- Creates a `dist/` directory.
- Generates `kex.json` (Index).
- Copies markdown files to `dist/`.

## `kex update`

Updates the Kex system documentation and agent configuration in an existing repository.

```bash
kex update [options]
```

- **System Docs (`contents/documentation/kex/*`)**: Updates to match the current binary version (Overwrite).
- **Agent Rules**: Updates based on the `.kex.yaml` strategies (Overwrite).
- See `.kex.yaml` configuration to customize behavior.
