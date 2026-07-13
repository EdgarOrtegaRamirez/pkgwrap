package pkgman

import (
	"strings"
)

// PipManager wraps the Python pip package manager.
type PipManager struct{}

func (p *PipManager) Name() string        { return "pip" }
func (p *PipManager) Description() string { return "Python package installer" }

func (p *PipManager) Available() bool {
	return commandExists("pip3") || commandExists("pip")
}

func (p *PipManager) pipBin() string {
	if commandExists("pip3") {
		return "pip3"
	}
	return "pip"
}

func (p *PipManager) Install(pkg string) error {
	return runCmd(p.pipBin(), "install", pkg)
}

func (p *PipManager) Search(query string) ([]SearchResult, error) {
	out, err := runCmdOutput(p.pipBin(), "search", query)
	if err != nil {
		// pip search may be disabled in newer versions
		// fallback to searching via pip index
		out2, err2 := runCmdOutput(p.pipBin(), "index", "versions", query)
		if err2 != nil {
			return nil, err
		}
		out = out2
	}
	results := []SearchResult{}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Parse "package_name (version)  - Description"
		parts := strings.SplitN(line, " - ", 2)
		nameVer := parts[0]
		desc := ""
		if len(parts) > 1 {
			desc = strings.TrimSpace(parts[1])
		}
		// Extract name and version from "name (version)"
		name := strings.TrimSpace(nameVer)
		version := ""
		if idx := strings.Index(nameVer, " ("); idx > 0 {
			name = strings.TrimSpace(nameVer[:idx])
			version = strings.TrimRight(strings.TrimSpace(nameVer[idx+2:]), ")")
		}
		results = append(results, SearchResult{
			Name:        name,
			Version:     version,
			Description: desc,
			Manager:     "pip",
		})
	}
	return results, nil
}

func (p *PipManager) Info(pkg string) (*PackageInfo, error) {
	out, err := runCmdOutput(p.pipBin(), "show", pkg)
	if err != nil {
		return nil, err
	}
	return parsePipShowOutput(out, pkg), nil
}

func (p *PipManager) ListInstalled() ([]InstalledPackage, error) {
	out, err := runCmdOutput(p.pipBin(), "list", "--format=columns")
	if err != nil {
		return nil, err
	}
	results := []InstalledPackage{}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i, line := range lines {
		if i == 0 || i == 1 {
			continue // skip header and separator
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			results = append(results, InstalledPackage{
				Name:    fields[0],
				Version: fields[1],
			})
		}
	}
	return results, nil
}

func (p *PipManager) Update() error {
	// pip doesn't have an index update; instead upgrade all outdated packages
	out, err := runCmdOutput(p.pipBin(), "list", "--outdated", "--format=columns")
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines[2:] { // skip header
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			runCmd(p.pipBin(), "install", "--upgrade", fields[0])
		}
	}
	return nil
}

func (p *PipManager) Remove(pkg string) error {
	return runCmd(p.pipBin(), "uninstall", "-y", pkg)
}