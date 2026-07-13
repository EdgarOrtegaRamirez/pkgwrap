package pkgman

// GoManager wraps the Go install command.
type GoManager struct{}

func (g *GoManager) Name() string        { return "go" }
func (g *GoManager) Description() string { return "Go package installer (go install)" }

func (g *GoManager) Available() bool {
	return commandExists("go")
}

func (g *GoManager) Install(pkg string) error {
	return runCmd("go", "install", pkg+"@latest")
}

func (g *GoManager) Search(query string) ([]SearchResult, error) {
	// Go doesn't have a built-in search; use pkg.go.dev via simple HTTP
	return nil, nil
}

func (g *GoManager) Info(pkg string) (*PackageInfo, error) {
	// Go doesn't have a built-in info command
	return nil, nil
}

func (g *GoManager) ListInstalled() ([]InstalledPackage, error) {
	// Go installs binaries to GOPATH/bin; no easy list command
	return nil, nil
}

func (g *GoManager) Update() error {
	return nil // no update command
}

func (g *GoManager) Remove(pkg string) error {
	// Go doesn't have a remove command
	return nil
}