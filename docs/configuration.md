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

```yaml
update:
  strategies:
    kex: all
    antigravity: coding-only
    cursor: all
```

- **strategies**: A map of **Directives** to update scopes.
    - **Keys**:
        - `kex`: Manages Kex system documentation.
        - `antigravity`, `cursor`, `claude`: Manages agent-specific rule files.
    - **Values (Directives)**:
        - `all`: Enable all categories (Coding + Documentation).
        - `coding-only`: Enable only Coding rules.
        - `documentation-only`: Enable only Documentation rules.
        - `none` (or omitted): Do not generate/update files for this key.

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


