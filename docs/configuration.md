# Configuration

## `.kex.yaml`

The configuration file `kex.yaml` (or `.kex.yaml`) lives at the root of your project.

```yaml
```yaml
sources:
  - contents
  - docs/guidelines
baseURL: https://example.com/docs/
```

### `sources` (Optional)

Specifies the directories containing your markdown documentation files.

- **Type**: `[]string` (List of strings)
- **Default**: `["contents"]`
- **Description**: Kex will recursively index all `.md` files within these directories. These paths are relative to the directory where `.kex.yaml` resides.
- **Description**: Kex will recursively index all `.md` files within these directories. These paths are relative to the directory where `.kex.yaml` resides.

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

```yaml
update:
  documents:
    kex: all
  ai-mcp-rules:
    targets: [antigravity, claude]
    scopes: [coding, documentation]
  ai-skills:
    targets: [claude]
    keywords: [go, typescript]
```

- **documents**:
    - `kex`: Manages Kex system documentation (`all` or `none`).
- **ai-mcp-rules**:
    - `targets`: List of agents to generate static rules for (e.g., `[antigravity, claude]`).
    - `scopes`: List of rule scopes to enforce (e.g., `[coding, documentation]`).
- **ai-skills**:
    - `targets`: List of agents to generate dynamic skills for (e.g., `[claude]`).
    - `keywords`: List of keywords to filter skills by.

> **Note**: `kex init` automatically generates this configuration based on your selections.

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


