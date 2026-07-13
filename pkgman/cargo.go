package pkgman

import "strings"

// CargoManager wraps the Rust cargo package manager.
type CargoManager struct{}

func (c *CargoManager) Name() string        { return "cargo" }
func (c *CargoManager) Description() string { return "Rust package manager" }

func (c *CargoManager) Available() bool {
	return commandExists("cargo")
}

func (c *CargoManager) Install(pkg string) error {
	return runCmd("cargo", "install", pkg)
}

func (c *CargoManager) Search(query string) ([]SearchResult, error) {
	out, err := runCmdOutput("cargo", "search", query, "--limit", "20")
	if err != nil {
		return nil, err
	}
	return parseCargoSearchOutput(out), nil
}

func (c *CargoManager) Info(pkg string) (*PackageInfo, error) {
	out, err := runCmdOutput("cargo", "search", pkg, "--limit", "5")
	if err != nil {
		return nil, err
	}
	return parseCargoSearchInfo(out, pkg), nil
}

func (c *CargoManager) ListInstalled() ([]InstalledPackage, error) {
	out, err := runCmdOutput("cargo", "install", "--list")
	if err != nil {
		return nil, err
	}
	results := []InstalledPackage{}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, " ") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			name := parts[0]
			version := ""
			if len(parts) >= 2 {
				version = strings.TrimRight(parts[1], ":")
			}
			results = append(results, InstalledPackage{
				Name:    name,
				Version: version,
			})
		}
	}
	return results, nil
}

func (c *CargoManager) Update() error {
	// cargo update upgrades dependencies in the current project
	return nil
}

func (c *CargoManager) Remove(pkg string) error {
	return runCmd("cargo", "uninstall", pkg)
}
