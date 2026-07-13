# pkgwrap 🔀

[![CI](https://github.com/EdgarOrtegaRamirez/pkgwrap/actions/workflows/ci.yml/badge.svg)](https://github.com/EdgarOrtegaRamirez/pkgwrap/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/EdgarOrtegaRamirez/pkgwrap)](https://goreportcard.com/report/github.com/EdgarOrtegaRamirez/pkgwrap)

**Universal package manager wrapper** — install, search, and manage packages across apt, brew, pip, npm, cargo, and go install from a single CLI.

```bash
pkgwrap install ripgrep          # Auto-detect best manager
pkgwrap search "json parser"
pkgwrap info curl
pkgwrap list --manager pip       # List pip packages
pkgwrap update                   # Update all package managers
pkgwrap managers                 # List available managers
```

## Features

- **🔍 Auto-detect** — Automatically finds the best package manager for each package
- **📦 6 managers** — apt, brew, pip, npm, cargo, go install
- **🔎 Search** — Search for packages across all managers or a specific one
- **📋 List** — View installed packages from any manager
- **ℹ️ Info** — Detailed package information (description, version, license, homepage)
- **📤 JSON output** — Machine-readable output for scripting and CI
- **🚦 CI-ready** — Non-zero exit codes for automation
- **⚡ No runtime deps** — Uses only Go standard library + cobra for CLI

## Installation

```bash
# From source
go install github.com/EdgarOrtegaRamirez/pkgwrap@latest

# Or build from repo
git clone https://github.com/EdgarOrtegaRamirez/pkgwrap.git
cd pkgwrap
go build -o pkgwrap .
sudo mv pkgwrap /usr/local/bin/
```

## Usage

### Install packages

```bash
# Auto-detect manager
pkgwrap install ripgrep

# Force a specific manager
pkgwrap install eslint --manager npm
pkgwrap install -m pip requests numpy
```

### Search for packages

```bash
# Search across all managers
pkgwrap search "json parser"

# Search in a specific manager
pkgwrap search lint --manager npm
pkgwrap search -m pip django
```

### Get package info

```bash
pkgwrap info curl
pkgwrap info typescript --manager npm
```

### List installed packages

```bash
pkgwrap list
pkgwrap list --manager pip
pkgwrap list -m npm --json
```

### Update

```bash
pkgwrap update                    # Update all
pkgwrap update --manager apt     # Update apt only
```

### Remove packages

```bash
pkgwrap remove ripgrep
pkgwrap remove -m pip requests
```

### List available managers

```bash
pkgwrap managers
pkgwrap managers --json
```

## Supported Managers

| Manager | Name      | Status |
|---------|-----------|--------|
| apt     | apt-get   | ✓      |
| brew    | Homebrew  | ✓      |
| pip     | Python    | ✓      |
| npm     | Node.js   | ✓      |
| cargo   | Rust      | ✓      |
| go      | Go        | ✓      |

## Development

```bash
go build -o pkgwrap .
go test ./... -v
go vet ./...
```

## License

MIT