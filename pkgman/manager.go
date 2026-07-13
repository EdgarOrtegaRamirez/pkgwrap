package pkgman

// PackageInfo represents detailed information about a package.
type PackageInfo struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	Manager          string `json:"manager"`
	Description      string `json:"description,omitempty"`
	Homepage         string `json:"homepage,omitempty"`
	License          string `json:"license,omitempty"`
	InstalledVersion string `json:"installed_version,omitempty"`
}

// SearchResult represents a package search result.
type SearchResult struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Manager     string `json:"manager"`
}

// InstalledPackage represents an installed package.
type InstalledPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Manager defines the interface for a package manager wrapper.
type Manager interface {
	// Name returns the package manager name (e.g., "apt", "pip", "npm").
	Name() string

	// Description returns a short description of the package manager.
	Description() string

	// Available returns true if this package manager is installed and usable.
	Available() bool

	// Install installs the specified package(s).
	Install(pkg string) error

	// Search searches for packages matching the query.
	Search(query string) ([]SearchResult, error)

	// Info returns detailed information about a package.
	Info(pkg string) (*PackageInfo, error)

	// ListInstalled returns all installed packages from this manager.
	ListInstalled() ([]InstalledPackage, error)

	// Update updates package indexes or upgrades packages.
	Update() error

	// Remove removes a package.
	Remove(pkg string) error
}

// Detect returns all available package managers on this system.
func Detect() []Manager {
	return []Manager{
		&AptManager{},
		&BrewManager{},
		&PipManager{},
		&NpmManager{},
		&CargoManager{},
		&GoManager{},
	}
}
