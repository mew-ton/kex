# Installation

You can install Kex using one of the following methods.

## Binary Download

You can download the latest pre-compiled binaries from the [GitHub Releases](https://github.com/mew-ton/kex/releases) page.

1. Download the archive matching your OS and Architecture (e.g., `_linux_amd64.tar.gz`, `_darwin_arm64.tar.gz`).
2. Extract the archive.
3. Move the `kex` binary to a directory in your `PATH` (e.g., `/usr/local/bin`).

```bash
# Example for Linux amd64
tar -xzf kex_linux_amd64.tar.gz
sudo mv kex /usr/local/bin/
```

## Go Install

If you have Go installed (1.25+), you can install the latest version directly:

```bash
go install github.com/mew-ton/kex/cmd/kex@latest
```

Ensure `$(go env GOPATH)/bin` is in your `PATH`.
