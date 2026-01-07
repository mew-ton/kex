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
    "contents/documentation/kex/*": "overwrite"
    "AGENTS.md": "marker-update"
```

- **strategies**: A map of glob patterns to update strategies.
  - `overwrite`: Replaces the file with the template.
  - `marker-update`: Updates content between markers (Designed for `AGENTS.md`).
  - `append`: Appends new content to the end if missing.
  - `skip`: Does nothing.

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


