package pkgman

// AptManager wraps the apt package manager.
type AptManager struct{}

func (a *AptManager) Name() string        { return "apt" }
func (a *AptManager) Description() string { return "Debian/Ubuntu package manager" }

func (a *AptManager) Available() bool {
	return commandExists("apt-get")
}

func (a *AptManager) Install(pkg string) error {
	return runCmd("apt-get", "install", "-y", pkg)
}

func (a *AptManager) Search(query string) ([]SearchResult, error) {
	out, err := runCmdOutput("apt-cache", "search", query)
	if err != nil {
		return nil, err
	}
	return parseAptSearchOutput(out), nil
}

func (a *AptManager) Info(pkg string) (*PackageInfo, error) {
	out, err := runCmdOutput("apt-cache", "show", pkg)
	if err != nil {
		return nil, err
	}
	return parseAptShowOutput(out, pkg), nil
}

func (a *AptManager) ListInstalled() ([]InstalledPackage, error) {
	out, err := runCmdOutput("dpkg-query", "-W", "-f=${Package} ${Version}\n")
	if err != nil {
		return nil, err
	}
	return parseDpkgListOutput(out), nil
}

func (a *AptManager) Update() error {
	return runCmd("apt-get", "update")
}

func (a *AptManager) Remove(pkg string) error {
	return runCmd("apt-get", "remove", "-y", pkg)
}
