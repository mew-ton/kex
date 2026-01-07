# Use Case: Monorepo / Local

This is the simplest setup. Documentation lives alongside the code it governs.

## Setup

1. Run `kex init` in your project root.
   - This creates `.kex.yaml` and `AGENTS.md`.
2. Write docs in `contents/`.
3. Configure your editor to run `kex start .` (or just `kex start`).

## Workflow

### 1. Setup Agent Instructions (`AGENTS.md`)

The generated `AGENTS.md` serves two purposes:
1.  **Docs Maintenance**: Instructions on how to add/edit files in `contents/`.
2.  **Coding Guidelines**: Prompts to be added to your agent configuration (e.g., `.cursorrules`) so it knows how to search for guidelines while coding.

**Action**: Copy the relevant sections from `AGENTS.md` into your agent's system prompt or configuration.

### 2. Daily Workflow

1.  **Write Code**: AI Agent queries `kex` for guidelines based on the prompt you added.
2.  **Update Docs**: If a rule changes, update the markdown file in the same PR.
3.  **CI/CD**: Run `kex check` to ensure documentation integrity.

### 3. CI Check Example

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

## Pros/Cons

- **Pros**: Versioned with code, zero latency, simple setup.
- **Cons**: Difficult to share rules across multiple repositories.
