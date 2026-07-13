package pkgman

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// runCmd runs a command with arguments, showing output to stdout/stderr.
func runCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// runCmdOutput runs a command and returns its combined output.
func runCmdOutput(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// parseAptSearchOutput parses apt-cache search output.
func parseAptSearchOutput(out string) []SearchResult {
	results := []SearchResult{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: "package - description"
		parts := strings.SplitN(line, " - ", 2)
		name := strings.TrimSpace(parts[0])
		desc := ""
		if len(parts) > 1 {
			desc = strings.TrimSpace(parts[1])
		}
		results = append(results, SearchResult{
			Name:        name,
			Description: desc,
			Manager:     "apt",
		})
	}
	return results
}

// parseAptShowOutput parses apt-cache show output.
func parseAptShowOutput(out string, pkg string) *PackageInfo {
	info := &PackageInfo{
		Name:    pkg,
		Manager: "apt",
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Version:") {
			info.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
		} else if strings.HasPrefix(line, "Description:") {
			info.Description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
		} else if strings.HasPrefix(line, "Homepage:") {
			info.Homepage = strings.TrimSpace(strings.TrimPrefix(line, "Homepage:"))
		}
	}
	return info
}

// parseDpkgListOutput parses dpkg-query list output.
func parseDpkgListOutput(out string) []InstalledPackage {
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
	return results
}

// parseBrewInfoJSON parses brew info --json output.
func parseBrewInfoJSON(out string, pkg string) *PackageInfo {
	info := &PackageInfo{
		Name:    pkg,
		Manager: "brew",
	}

	var brewInfo struct {
		Name        string `json:"name"`
		Versions    struct {
			Stable string `json:"stable"`
		} `json:"versions"`
		Description string `json:"desc"`
		Homepage    string `json:"homepage"`
		License     string `json:"license"`
	}
	if err := json.Unmarshal([]byte(out), &brewInfo); err == nil {
		info.Name = brewInfo.Name
		info.Version = brewInfo.Versions.Stable
		info.Description = brewInfo.Description
		info.Homepage = brewInfo.Homepage
		info.License = brewInfo.License
	}

	// Try as array
	var brewInfoArr []struct {
		Name        string `json:"name"`
		Versions    struct {
			Stable string `json:"stable"`
		} `json:"versions"`
		Description string `json:"desc"`
		Homepage    string `json:"homepage"`
		License     string `json:"license"`
	}
	if err := json.Unmarshal([]byte(out), &brewInfoArr); err == nil && len(brewInfoArr) > 0 {
		info.Name = brewInfoArr[0].Name
		info.Version = brewInfoArr[0].Versions.Stable
		info.Description = brewInfoArr[0].Description
		info.Homepage = brewInfoArr[0].Homepage
		info.License = brewInfoArr[0].License
	}

	return info
}

// parsePipShowOutput parses pip show output.
func parsePipShowOutput(out string, pkg string) *PackageInfo {
	info := &PackageInfo{
		Name:    pkg,
		Manager: "pip",
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name:") {
			info.Name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
		} else if strings.HasPrefix(line, "Version:") {
			info.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
		} else if strings.HasPrefix(line, "Summary:") {
			info.Description = strings.TrimSpace(strings.TrimPrefix(line, "Summary:"))
		} else if strings.HasPrefix(line, "Home-page:") {
			info.Homepage = strings.TrimSpace(strings.TrimPrefix(line, "Home-page:"))
		} else if strings.HasPrefix(line, "License:") {
			info.License = strings.TrimSpace(strings.TrimPrefix(line, "License:"))
		}
	}
	return info
}

// parseCargoSearchOutput parses cargo search output.
func parseCargoSearchOutput(out string) []SearchResult {
	results := []SearchResult{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "    ") {
			continue
		}
		// Format: "name = \"version\" #description"
		if idx := strings.Index(line, " = "); idx > 0 {
			name := strings.TrimSpace(line[:idx])
			rest := strings.TrimSpace(line[idx+3:])
			version := ""
			desc := ""
			if strings.HasPrefix(rest, "\"") {
				parts := strings.SplitN(rest[1:], "\"", 2)
				if len(parts) >= 1 {
					version = parts[0]
				}
				if len(parts) >= 2 {
					desc = strings.TrimPrefix(parts[1], " #")
				}
			}
			results = append(results, SearchResult{
				Name:        name,
				Version:     version,
				Description: desc,
				Manager:     "cargo",
			})
		}
	}
	return results
}

// parseCargoSearchInfo parses cargo search output for a single package.
func parseCargoSearchInfo(out, pkg string) *PackageInfo {
	results := parseCargoSearchOutput(out)
	for _, r := range results {
		if r.Name == pkg {
			return &PackageInfo{
				Name:        r.Name,
				Version:     r.Version,
				Description: r.Description,
				Manager:     "cargo",
			}
		}
	}
	return nil
}