# Configuration

## `.kex.yaml`

The configuration file `kex.yaml` (or `.kex.yaml`) lives at the root of your project.

```yaml
root: contents
baseURL: https://example.com/docs/
```

### `root` (Optional)

Specifies the directory containing your markdown documentation files.

- **Type**: `string`
- **Default**: `contents`
- **Description**: Kex will recursively index all `.md` files within this directory. This path is relative to the directory where `kex start` or `kex check` is run (usually the project root).

### `baseURL` (Optional)

Defines the base URL for the remote hosting location of your documentation.

- **Type**: `string`
- **Example**: `https://my-org.github.io/my-repo/`
- **Description**:
    - Used by `kex generate` to create absolute URLs in the `kex.json` index.
    - If omitted, `kex generate` will produce relative paths (which generally works fine for most static site setups).
    - Necessary if you plan to consume the `kex.json` from a different domain or context where relative paths might break.

### `remoteToken` (Optional)

Authentication token for private repositories (e.g. GitHub Private Pages).

- **Type**: `string`
- **Description**: If set, this token will be sent as a Bearer token in the `Authorization` header when fetching remote documents.

### `update` (Optional)

Configures the behavior of `kex update`.

### `update` (Optional)

Configures the behavior of `kex update`.

```yaml
update:
  strategies:
    ".agent/rules/kex-coding.md": "skip"
    ".claude/rules/kex-coding.md": "skip"
    "contents/documentation/kex/choose-effective-keywords.md": "overwrite"
```

- **strategies**: A map of **Canonical Paths** to update strategies.
    - **Rule Files**: Use the specific path for the agent type (e.g., `.agent/rules/kex-coding.md` for General, `.claude/rules/kex-coding.md` for Claude/Cursor).
    - **System Documentation**: Always use the path starting with `contents/` (e.g., `contents/documentation/kex/choose-effective-keywords.md`), even if you have configured a custom `root`.
    - **Strategies**:
        - `overwrite`: Replaces the file with the template.
        - `skip`: Creates the file if missing, but strictly preserves existing content if the file exists.
    - **Defaults**:
        - **Rule Files**: Default to `skip` to preserve user customizations.
        - **System Documentation**: Defaults to `overwrite`.


### `agent` (Optional)

Configures the AI Agent guidelines generation.

```yaml
agent:
  type: general # or "claude"
  scopes:
    - coding
    - documentation
```

- **type**: The type of agent guidelines to generate.
    - `general`: Standard `AGENTS.md`.
    - `claude`: `CLAUDE.md` tailored for Anthropic Claude.
- **scopes**: A list of guideline sections to include.
    - `coding`: Design and Implementation phase guidelines.
    - `documentation`: Documentation phase rules.

### `logging` (Optional)

Configures server logging.

```yaml
logging:
  file: ".kex/kex.log"
  level: "info"
```

- **file**: Path to the log file. If omitted, logs are written to Stderr.
- **level**: Log level (default: `info`).


## Environment Variables

Kex supports the following environment variables:

### `KEX_REMOTE_TOKEN`

- **Purpose**: Authenticate against private remote documentation endpoints.
- **Priority**: Takes precedence over `.kex.yaml`'s `remoteToken`.
- **Usage**:
  ```bash
  export KEX_REMOTE_TOKEN="your-secret-token"
  kex start https://private.example.com/docs/
  ```


