# pkgwrap — AGENTS.md

## Project Overview

pkgwrap is a universal package manager wrapper CLI in Go. It unifies apt, brew, pip, npm, cargo, and go install under a single interface.

## Architecture

```
pkgwrap/
├── main.go                  # Entry point
├── cmd/
│   └── root.go              # Cobra commands (install, search, info, list, update, remove, managers)
├── pkgman/
│   ├── manager.go           # Manager interface + Detect()
│   ├── apt.go               # apt-get implementation
│   ├── brew.go              # Homebrew implementation
│   ├── pip.go               # pip implementation
│   ├── npm.go               # npm implementation
│   ├── npm_parse.go         # npm JSON/text parsing helpers
│   ├── cargo.go             # Cargo implementation
│   ├── go.go                # go install implementation
│   ├── utils.go             # Shared utilities (commandExists, runCmd, parser functions)
│   ├── pkgman_test.go       # 20 unit tests
├── go.mod / go.sum
├── README.md
├── LICENSE
└── .github/workflows/ci.yml
```

## Key Design Decisions

1. **Interface-based** — `Manager` interface allows adding new package managers easily
2. **Auto-detect** — `Detect()` checks PATH for each manager binary
3. **Text parsing** — Parses CLI output from each manager rather than using libraries (avoids dependency hell)
4. **Fallback chain** — `install` and `remove` try managers in order until one succeeds
5. **JSON output** — Every command supports `--json` for CI/CD pipelines
6. **Zero runtime dependencies** — Only `github.com/spf13/cobra` for CLI (stdlib for everything else)

## Dependencies

- `github.com/spf13/cobra` — CLI framework

## Build & Test

```bash
go build -o pkgwrap .
go test ./... -v
go vet ./...
```

## Adding a New Manager

1. Create `pkgman/<name>.go` implementing the `Manager` interface
2. Add to `Detect()` in `manager.go`
3. Add tests in `pkgman/pkgman_test.go`
4. Run `go test ./... -v`
5. Update README.md with the new manager

## Common Tasks

### Fix CLI parsing for a new manager
Each manager's `Search()`, `Info()`, and `ListInstalled()` methods parse text output. If output format changes, update the corresponding `parse*` function in utils.go or the manager's own file.

### Add a new command
1. Add a `cobra.Command` variable in `cmd/root.go`
2. Register it in `init()` with `rootCmd.AddCommand()`
3. Implement the handler using the Manager interface
4. Add to `--json` support if applicable