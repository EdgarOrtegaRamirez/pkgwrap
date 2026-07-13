package pkgman

import (
	"encoding/json"
	"fmt"
	"strings"
)

// NpmSearchResult represents a single npm search result from JSON output.
type npmPackageJSON struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func parseNpmSearchJSON(data string) ([]SearchResult, error) {
	var pkgs []npmPackageJSON
	if err := json.Unmarshal([]byte(data), &pkgs); err != nil {
		return nil, err
	}
	results := make([]SearchResult, len(pkgs))
	for i, p := range pkgs {
		results[i] = SearchResult{
			Name:        p.Name,
			Version:     p.Version,
			Description: p.Description,
			Manager:     "npm",
		}
	}
	return results, nil
}

func parseNpmSearchText(data string) []SearchResult {
	results := []SearchResult{}
	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			name := fields[0]
			desc := strings.Join(fields[1:], " ")
			results = append(results, SearchResult{
				Name:        name,
				Description: desc,
				Manager:     "npm",
			})
		}
	}
	return results
}

func parseNpmViewJSON(data string, pkg string) (*PackageInfo, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	info := &PackageInfo{
		Name:    pkg,
		Manager: "npm",
	}
	if v, ok := raw["version"].(string); ok {
		info.Version = v
	}
	if v, ok := raw["description"].(string); ok {
		info.Description = v
	}
	if v, ok := raw["homepage"].(string); ok {
		info.Homepage = v
	}
	if v, ok := raw["license"]; ok {
		if s, ok := v.(string); ok {
			info.License = s
		}
	}
	return info, nil
}

func parseNpmListJSON(data string) ([]InstalledPackage, error) {
	var root struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(data), &root); err != nil {
		return nil, err
	}
	results := []InstalledPackage{}
	for name, dep := range root.Dependencies {
		results = append(results, InstalledPackage{
			Name:    name,
			Version: dep.Version,
		})
	}
	return results, nil
}

func parseAndUpgradeNpmOutdated(data string) error {
	var outdated map[string]struct {
		Current string `json:"current"`
		Wanted  string `json:"wanted"`
		Latest  string `json:"latest"`
	}
	if err := json.Unmarshal([]byte(data), &outdated); err != nil {
		return nil
	}
	for pkg := range outdated {
		if err := runCmd("npm", "install", "-g", pkg); err != nil {
			return fmt.Errorf("failed to upgrade %s: %v", pkg, err)
		}
	}
	return nil
}