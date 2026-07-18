package pkgman

import "strings"

// BrewManager wraps the Homebrew package manager.
type BrewManager struct{}

func (b *BrewManager) Name() string        { return "brew" }
func (b *BrewManager) Description() string { return "Homebrew package manager (macOS/Linux)" }

func (b *BrewManager) Available() bool {
	return commandExists("brew")
}

func (b *BrewManager) Install(pkg string) error {
	return runCmd("brew", "install", pkg)
}

func (b *BrewManager) Search(query string) ([]SearchResult, error) {
	out, err := runCmdOutput("brew", "search", query)
	if err != nil {
		return nil, err
	}
	results := []SearchResult{}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "==> Formulae" || line == "==> Casks" || strings.HasPrefix(line, "==>") {
			continue
		}
		// brew search output is one package per line
		results = append(results, SearchResult{
			Name:    strings.TrimSpace(line),
			Manager: "brew",
		})
	}
	return results, nil
}

func (b *BrewManager) Info(pkg string) (*PackageInfo, error) {
	out, err := runCmdOutput("brew", "info", "--json=v2", pkg)
	if err != nil {
		return nil, err
	}
	return parseBrewInfoJSON(out, pkg), nil
}

func (b *BrewManager) ListInstalled() ([]InstalledPackage, error) {
	out, err := runCmdOutput("brew", "list", "--versions")
	if err != nil {
		return nil, err
	}
	results := []InstalledPackage{}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			results = append(results, InstalledPackage{
				Name:    parts[0],
				Version: parts[1],
			})
		}
	}
	return results, nil
}

func (b *BrewManager) Update() error {
	return runCmd("brew", "update")
}

func (b *BrewManager) Remove(pkg string) error {
	return runCmd("brew", "uninstall", pkg)
}
