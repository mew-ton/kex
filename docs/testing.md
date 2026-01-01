# Testing Guidelines

This document outlines the testing strategy and guidelines for the `kex` project.

## Philosophy: Test First & Maintenance Priority

The reliability of the `kex` CLI is critical as it is used by the Document Librarian MCP. To ensure a stable "vibe coding" experience:

1.  **Test First**: Tests **must** pass and the project **must** build at all times. The MCP relies on the latest CLI version, so a broken main branch directly impacts usability.
2.  **Maintenance Priority**: Maintaining existing tests is as important as writing new code. If a change breaks a test, fixing the test (or code) is the priority.

## E2E Testing

End-to-End (E2E) tests are located in the `e2e/` directory. They build the `kex` binary and run it against fixtures.

### Running Tests

Use the `Makefile` target to run E2E tests:

```bash
make e2e
```

### Fixtures

-   **`examples/`**: Use the content in `examples/` for general, positive test cases. This ensures our examples are always working.
-   **`e2e/fixtures/`**: Use this directory for specific test cases that require isolation, negative testing (e.g., invalid configs), or scenarios that shouldn't clutter the main examples. Create a subdirectory for each test case (e.g., `e2e/fixtures/invalid-config`).

## Best Practices

-   **Avoid Flakiness**: Prefer checking **exit codes** and **side effects** (e.g., file creation) over strict text matching on `stdout`. Text output often changes and can lead to brittle tests.
-   **Future Helper**: For complex commands like `check`, we plan to implement a JSON output mode to allow precise, structural assertions closer to unit tests.


## Pre-commit Checks (Lefthook)

We use [Lefthook](https://github.com/evilmartians/lefthook) to enforce quality before every commit.
The pre-commit hook runs:

-   `go vet ./...` (Linting)
-   `go test -v ./...` (Unit Tests)
-   `go build` (Build verification)

To install hooks (if you haven't already):

```bash
make init
```
