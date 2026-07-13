package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/EdgarOrtegaRamirez/pkgwrap/pkgman"
	"github.com/spf13/cobra"
)

var (
	managerFlag string
	jsonOutput  bool
	quietOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "pkgwrap",
	Short: "Universal package manager wrapper",
	Long: `pkgwrap is a universal package manager wrapper that unifies
multiple package managers (apt, brew, pip, npm, cargo, go install, etc.)
under a single CLI interface.

Examples:
  pkgwrap install nodejs          # Auto-detect best manager
  pkgwrap install eslint --manager npm
  pkgwrap search "json parser"
  pkgwrap info curl
  pkgwrap managers                # List available managers
  pkgwrap list --manager pip      # List packages from a specific manager
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&managerFlag, "manager", "m", "", "Specific package manager to use (auto-detect if empty)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&quietOutput, "quiet", "q", false, "Suppress non-error output")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(managersCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(removeCmd)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func getManagers(selected string) []pkgman.Manager {
	all := pkgman.Detect()
	if selected == "" {
		return all
	}
	for _, m := range all {
		if strings.EqualFold(m.Name(), selected) {
			return []pkgman.Manager{m}
		}
	}
	fmt.Fprintf(os.Stderr, "Package manager %q not found (available: %s)\n",
		selected, strings.Join(managerNames(all), ", "))
	os.Exit(1)
	return nil
}

func managerNames(managers []pkgman.Manager) []string {
	names := make([]string, len(managers))
	for i, m := range managers {
		names[i] = m.Name()
	}
	return names
}

// --- install ---

var installCmd = &cobra.Command{
	Use:   "install <package> [package...]",
	Short: "Install packages",
	Long: `Install one or more packages using the best available package manager.
Use --manager to force a specific manager.

Examples:
  pkgwrap install ripgrep
  pkgwrap install eslint prettier --manager npm
  pkgwrap install -m pip requests
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		managers := getManagers(managerFlag)
		exitCode := 0
		for _, pkg := range args {
			installed := false
			for _, m := range managers {
				if !m.Available() {
					continue
				}
				if !quietOutput {
					fmt.Fprintf(os.Stderr, "[%s] Installing %s...\n", m.Name(), pkg)
				}
				if err := m.Install(pkg); err != nil {
					fmt.Fprintf(os.Stderr, "[%s] Failed to install %s: %v\n", m.Name(), pkg, err)
					exitCode = 1
					continue
				}
				if !quietOutput {
					fmt.Printf("✓ %s installed via %s\n", pkg, m.Name())
				}
				installed = true
				break
			}
			if !installed {
				fmt.Fprintf(os.Stderr, "Could not install %s: no suitable package manager found\n", pkg)
				exitCode = 1
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

// --- search ---

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for packages across all managers",
	Long: `Search for packages matching the query across all available package managers.
Use --manager to search only a specific manager.

Examples:
  pkgwrap search "json parser"
  pkgwrap search lint --manager npm
  pkgwrap search -m pip django
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		managers := getManagers(managerFlag)
		exitCode := 0
		var results []pkgman.SearchResult

		for _, m := range managers {
			if !m.Available() {
				continue
			}
			res, err := m.Search(query)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] Search failed: %v\n", m.Name(), err)
				exitCode = 1
				continue
			}
			results = append(results, res...)
		}

		if jsonOutput {
			return printJSON(results)
		}

		if len(results) == 0 {
			fmt.Println("No results found.")
			return nil
		}

		// Group by manager
		byManager := make(map[string][]pkgman.SearchResult)
		for _, r := range results {
			byManager[r.Manager] = append(byManager[r.Manager], r)
		}

		managerNames := make([]string, 0, len(byManager))
		for name := range byManager {
			managerNames = append(managerNames, name)
		}
		sort.Strings(managerNames)

		for _, name := range managerNames {
			fmt.Printf("\n── %s ──\n", name)
			for _, r := range byManager[name] {
				desc := r.Description
				if len(desc) > 80 {
					desc = desc[:77] + "..."
				}
				fmt.Printf("  %-30s %s\n", r.Name, desc)
			}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

// --- info ---

var infoCmd = &cobra.Command{
	Use:   "info <package>",
	Short: "Show package information",
	Long: `Show detailed information about a package from the first available manager.
Use --manager to query a specific manager.

Examples:
  pkgwrap info curl
  pkgwrap info typescript --manager npm
  pkgwrap info -m pip requests
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		managers := getManagers(managerFlag)

		for _, m := range managers {
			if !m.Available() {
				continue
			}
			info, err := m.Info(pkg)
			if err != nil {
				continue
			}
			return printPackageInfo(info, jsonOutput)
		}

		return fmt.Errorf("package %q not found in any available package manager", pkg)
	},
}

// --- managers ---

var managersCmd = &cobra.Command{
	Use:   "managers",
	Short: "List available package managers",
	Long: `List all detected package managers and their availability on this system.

Examples:
  pkgwrap managers
  pkgwrap managers --json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		managers := pkgman.Detect()
		if jsonOutput {
			type managerInfo struct {
				Name        string `json:"name"`
				Available   bool   `json:"available"`
				Description string `json:"description"`
			}
			infos := make([]managerInfo, len(managers))
			for i, m := range managers {
				infos[i] = managerInfo{
					Name:        m.Name(),
					Available:   m.Available(),
					Description: m.Description(),
				}
			}
			return printJSON(infos)
		}

		fmt.Println("Available package managers:")
		fmt.Println()
		for _, m := range managers {
			status := "✓"
			if !m.Available() {
				status = "✗"
			}
			fmt.Printf("  %s %-12s %s\n", status, m.Name(), m.Description())
		}
		return nil
	},
}

// --- list ---

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Long: `List installed packages from all available package managers.
Use --manager to list packages from a specific manager.

Examples:
  pkgwrap list
  pkgwrap list --manager pip
  pkgwrap list -m npm --json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		managers := getManagers(managerFlag)
		exitCode := 0
		type pkgEntry struct {
			Manager string `json:"manager"`
			Package string `json:"package"`
			Version string `json:"version"`
		}
		var allPkgs []pkgEntry

		for _, m := range managers {
			if !m.Available() {
				continue
			}
			pkgs, err := m.ListInstalled()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] Failed to list packages: %v\n", m.Name(), err)
				exitCode = 1
				continue
			}
			for _, p := range pkgs {
				allPkgs = append(allPkgs, pkgEntry{
					Manager: m.Name(),
					Package: p.Name,
					Version: p.Version,
				})
			}
		}

		if jsonOutput {
			return printJSON(allPkgs)
		}

		if len(allPkgs) == 0 {
			fmt.Println("No packages found.")
			return nil
		}

		// Group by manager
		byManager := make(map[string][]pkgEntry)
		for _, p := range allPkgs {
			byManager[p.Manager] = append(byManager[p.Manager], p)
		}

		managerNames := make([]string, 0, len(byManager))
		for name := range byManager {
			managerNames = append(managerNames, name)
		}
		sort.Strings(managerNames)

		for _, name := range managerNames {
			fmt.Printf("\n── %s (%d packages) ──\n", name, len(byManager[name]))
			for _, p := range byManager[name] {
				fmt.Printf("  %-30s %s\n", p.Package, p.Version)
			}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

// --- update ---

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update package lists or upgrade packages",
	Long: `Update package manager indexes or upgrade all packages.
Use --manager to target a specific manager.

Examples:
  pkgwrap update                    # Update all managers
  pkgwrap update --manager apt     # Update apt only
  pkgwrap update --manager pip     # Upgrade all pip packages
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		managers := getManagers(managerFlag)
		exitCode := 0
		for _, m := range managers {
			if !m.Available() {
				continue
			}
			if !quietOutput {
				fmt.Fprintf(os.Stderr, "[%s] Updating...\n", m.Name())
			}
			if err := m.Update(); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] Update failed: %v\n", m.Name(), err)
				exitCode = 1
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

// --- remove ---

var removeCmd = &cobra.Command{
	Use:   "remove <package> [package...]",
	Short: "Remove packages",
	Long: `Remove one or more packages. Uses the same manager that was used to install.
Use --manager to force a specific manager.

Examples:
  pkgwrap remove ripgrep
  pkgwrap remove eslint --manager npm
  pkgwrap remove -m pip requests
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		managers := getManagers(managerFlag)
		exitCode := 0
		for _, pkg := range args {
			removed := false
			for _, m := range managers {
				if !m.Available() {
					continue
				}
				if !quietOutput {
					fmt.Fprintf(os.Stderr, "[%s] Removing %s...\n", m.Name(), pkg)
				}
				if err := m.Remove(pkg); err != nil {
					fmt.Fprintf(os.Stderr, "[%s] Failed to remove %s: %v\n", m.Name(), pkg, err)
					exitCode = 1
					continue
				}
				if !quietOutput {
					fmt.Printf("✓ %s removed via %s\n", pkg, m.Name())
				}
				removed = true
				break
			}
			if !removed {
				fmt.Fprintf(os.Stderr, "Could not remove %s: no suitable package manager found\n", pkg)
				exitCode = 1
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

func printPackageInfo(info *pkgman.PackageInfo, jsonFmt bool) error {
	if jsonFmt {
		return printJSON(info)
	}
	fmt.Printf("Name:        %s\n", info.Name)
	fmt.Printf("Version:     %s\n", info.Version)
	fmt.Printf("Manager:     %s\n", info.Manager)
	if info.Description != "" {
		fmt.Printf("Description: %s\n", info.Description)
	}
	if info.Homepage != "" {
		fmt.Printf("Homepage:    %s\n", info.Homepage)
	}
	if info.License != "" {
		fmt.Printf("License:     %s\n", info.License)
	}
	if info.InstalledVersion != "" {
		fmt.Printf("Installed:   %s\n", info.InstalledVersion)
	}
	return nil
}

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}