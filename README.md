# forklift

[![License](https://img.shields.io/github/license/homeport/forklift.svg)](https://github.com/homeport/forklift/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/forklift)](https://goreportcard.com/report/github.com/homeport/forklift)
[![Tests](https://github.com/homeport/forklift/workflows/Tests/badge.svg)](https://github.com/homeport/forklift/actions?query=workflow%3A%22Tests%22)
[![Codecov](https://img.shields.io/codecov/c/github/homeport/forklift/main.svg)](https://codecov.io/gh/homeport/forklift)
[![Go Reference](https://pkg.go.dev/badge/github.com/homeport/forklift.svg)](https://pkg.go.dev/github.com/homeport/forklift)
[![Release](https://img.shields.io/github/release/homeport/forklift.svg)](https://github.com/homeport/forklift/releases/latest)

## Description

Experimental tool to manipulate container images.

## Installation

### Homebrew

The `homeport/tap` has macOS and GNU/Linux pre-built binaries available:

```bash
brew install homeport/tap/forklift
```

### Pre-built binaries in GitHub

Prebuilt binaries can be [downloaded from the GitHub Releases section](https://github.com/homeport/forklift/releases/latest).

### Curl To Shell Convenience Script

There is a convenience script to download the latest release for Linux or macOS if you want to need it simple (you need `curl` and `jq` installed on your machine):

```bash
curl --silent --location https://raw.githubusercontent.com/homeport/forklift/main/hack/download.sh | bash
```

### Build from Source

You can install `forklift` from source using `go install`:

```bash
go install github.com/homeport/forklift/cmd/forklift@latest
```

_Please note:_ This will install `forklift` based on the latest available code base. Even though the goal is that the latest commit on the `main` branch should always be a stable and usable version, this is not the recommended way to install and use `forklift`. If you find an issue with this version, please make sure to note the commit SHA or date in the GitHub issue to indcate that it is not based on a released version. The version output will show `forklift version (development)` for `go install` based builds.

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get Go (`forklift` requires Go version 1.20 or greater)
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):
> The seven rules of a great Git commit message
>
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

### Running test cases and binaries generation

Run test cases:

```bash
ginkgo run ./...
```

Create binaries:

```bash
goreleaser build --clean --snapshot
```

## License

Licensed under [MIT License](https://github.com/homeport/forklift/blob/main/LICENSE)
