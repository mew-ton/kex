# Contributing to Kex

Thank you for your interest in contributing to Kex!

## Release Workflow

This project uses **Tagpr** for automated releases.

1.  **Pull Request**: When a PR is merged into `main`, a "Release PR" is automatically created/updated by the GitHub Actions bot. This PR increments the version in `README.md` (if applicable) and updates `CHANGELOG.md`.
2.  **Release**: When you merge the "Release PR", a new git tag is created, and **GoReleaser** automatically builds the binaries and publishes a GitHub Release.

## Design Documents

Please refer to [docs/design.md](docs/design.md) for the architectural design.
