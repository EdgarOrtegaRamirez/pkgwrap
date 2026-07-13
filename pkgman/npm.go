package pkgman

import "fmt"

// NpmManager wraps the npm package manager.
type NpmManager struct{}

func (n *NpmManager) Name() string        { return "npm" }
func (n *NpmManager) Description() string { return "Node.js package manager" }

func (n *NpmManager) Available() bool {
	return commandExists("npm")
}

func (n *NpmManager) Install(pkg string) error {
	return runCmd("npm", "install", "-g", pkg)
}

func (n *NpmManager) Search(query string) ([]SearchResult, error) {
	out, err := runCmdOutput("npm", "search", query, "--json")
	if err != nil {
		// npm search may not return JSON; try without --json
		out2, err2 := runCmdOutput("npm", "search", query)
		if err2 != nil {
			return nil, fmt.Errorf("npm search failed: %v", err)
		}
		return parseNpmSearchText(out2), nil
	}
	return parseNpmSearchJSON(out)
}

func (n *NpmManager) Info(pkg string) (*PackageInfo, error) {
	out, err := runCmdOutput("npm", "view", pkg, "--json")
	if err != nil {
		return nil, err
	}
	return parseNpmViewJSON(out, pkg)
}

func (n *NpmManager) ListInstalled() ([]InstalledPackage, error) {
	out, err := runCmdOutput("npm", "list", "-g", "--depth=0", "--json")
	if err != nil {
		return nil, err
	}
	return parseNpmListJSON(out)
}

func (n *NpmManager) Update() error {
	// npm update all global packages
	out, err := runCmdOutput("npm", "outdated", "-g", "--json")
	if err != nil {
		return nil // up to date or no packages
	}
	return parseAndUpgradeNpmOutdated(out)
}

func (n *NpmManager) Remove(pkg string) error {
	return runCmd("npm", "uninstall", "-g", pkg)
}
