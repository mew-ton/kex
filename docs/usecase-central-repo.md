# Use Case: Central Repository (GitHub Pages)

Ideal for organizations with shared guidelines across many services.

## Setup

1. Create a dedicated repository (e.g., `my-org/guidelines`).
2. Run `kex init`.
3. Configure GitHub Actions to run `kex generate dist` and deploy `dist/` to GitHub Pages.

### CI Check Example

Always run `kex check` on Pull Requests to ensure no invalid documents are merged.

Create `.github/workflows/check.yml`:

```yaml
name: Check Guidelines

on:
  pull_request:

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Kex
        run: go install github.com/mew-ton/kex/cmd/kex@latest

      - name: Check Docs
        run: kex check
```

### Deploy Workflow Example

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy Guidelines

on:
  push:
    branches: ["main"]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Kex
        run: go install github.com/mew-ton/kex/cmd/kex@latest

      - name: Generate Site
        run: kex generate dist

      - name: Setup Pages
        uses: actions/configure-pages@v5

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: './dist'

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Usage in Other Repos

In your consuming repositories, configure the MCP client to point to the remote URL:

```bash
kex start https://my-org.github.io/guidelines/
```

## Workflow

1. **Central Update**: Update rules in the guidelines repo.
2. **Deploy**: GitHub Actions publishes the new static site + `kex.json`.
3. **Consumption**: All connected agents immediately see the new rules.

## Pros/Cons

- **Pros**: Single source of truth, easy to update globally.
- **Cons**: Version drift (code might lag behind rules), network dependency.
