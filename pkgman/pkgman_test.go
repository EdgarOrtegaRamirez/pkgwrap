package pkgman

import (
	"testing"
)

func TestParseAptSearchOutput(t *testing.T) {
	input := "ripgrep - recursively searches directories for a regex pattern\nhtop - interactive process viewer\n"
	results := parseAptSearchOutput(input)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Name != "ripgrep" {
		t.Errorf("expected ripgrep, got %s", results[0].Name)
	}
	if results[0].Manager != "apt" {
		t.Errorf("expected apt manager, got %s", results[0].Manager)
	}
}

func TestParseAptShowOutput(t *testing.T) {
	input := "Package: curl\nVersion: 7.68.0\nDescription: command line tool for transferring data with URL syntax\nHomepage: https://curl.haxx.se\n"
	info := parseAptShowOutput(input, "curl")
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info.Name != "curl" {
		t.Errorf("expected curl, got %s", info.Name)
	}
	if info.Version != "7.68.0" {
		t.Errorf("expected 7.68.0, got %s", info.Version)
	}
	if info.Manager != "apt" {
		t.Errorf("expected apt manager, got %s", info.Manager)
	}
}

func TestParseDpkgListOutput(t *testing.T) {
	input := "curl 7.68.0-1ubuntu1\ngit 2.25.1\npython3 3.8.2\n"
	pkgs := parseDpkgListOutput(input)
	if len(pkgs) != 3 {
		t.Fatalf("expected 3 packages, got %d", len(pkgs))
	}
	if pkgs[1].Name != "git" {
		t.Errorf("expected git, got %s", pkgs[1].Name)
	}
	if pkgs[1].Version != "2.25.1" {
		t.Errorf("expected 2.25.1, got %s", pkgs[1].Version)
	}
}

func TestParseCargoSearchOutput(t *testing.T) {
	input := `serde = "1.0.217" # Serde is a framework for serializing and deserializing Rust data structures
tokio = "1.44.2" # An event-driven, non-blocking I/O platform for writing asynchronous I/O
`
	results := parseCargoSearchOutput(input)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Name != "serde" {
		t.Errorf("expected serde, got %s", results[0].Name)
	}
	if results[0].Manager != "cargo" {
		t.Errorf("expected cargo manager, got %s", results[0].Manager)
	}
	if results[1].Version != "1.44.2" {
		t.Errorf("expected 1.44.2, got %s", results[1].Version)
	}
}

func TestParsePipShowOutput(t *testing.T) {
	input := "Name: requests\nVersion: 2.25.1\nSummary: Python HTTP for Humans.\nHome-page: https://requests.readthedocs.io\nLicense: Apache 2.0\n"
	info := parsePipShowOutput(input, "requests")
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info.Name != "requests" {
		t.Errorf("expected requests, got %s", info.Name)
	}
	if info.Version != "2.25.1" {
		t.Errorf("expected 2.25.1, got %s", info.Version)
	}
	if info.Manager != "pip" {
		t.Errorf("expected pip manager, got %s", info.Manager)
	}
}

func TestParseNpmSearchJSON(t *testing.T) {
	input := `[{"name":"express","version":"4.18.2","description":"Fast, unopinionated, minimalist web framework"}]`
	results, err := parseNpmSearchJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "express" {
		t.Errorf("expected express, got %s", results[0].Name)
	}
	if results[0].Manager != "npm" {
		t.Errorf("expected npm manager, got %s", results[0].Manager)
	}
}

func TestParseNpmSearchText(t *testing.T) {
	input := "express  Fast, unopinionated, minimalist web framework\nlodash  The Lodash library exported as Node.js modules."
	results := parseNpmSearchText(input)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[1].Name != "lodash" {
		t.Errorf("expected lodash, got %s", results[1].Name)
	}
}

func TestParseBrewInfoJSON(t *testing.T) {
	input := `[{"name":"curl","versions":{"stable":"8.4.0"},"desc":"Command line tool for transferring data","homepage":"https://curl.se","license":"MIT"}]`
	info := parseBrewInfoJSON(input, "curl")
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info.Name != "curl" {
		t.Errorf("expected curl, got %s", info.Name)
	}
	if info.Version != "8.4.0" {
		t.Errorf("expected 8.4.0, got %s", info.Version)
	}
}

func TestCommandExists(t *testing.T) {
	if !commandExists("go") {
		t.Error("expected go to exist")
	}
	if commandExists("nonexistent-command-xyz123") {
		t.Error("expected nonexistent command to not exist")
	}
}

func TestManagerDetect(t *testing.T) {
	managers := Detect()
	expectedNames := []string{"apt", "brew", "pip", "npm", "cargo", "go"}
	if len(managers) != len(expectedNames) {
		t.Fatalf("expected %d managers, got %d", len(expectedNames), len(managers))
	}
	for i, m := range managers {
		if m.Name() != expectedNames[i] {
			t.Errorf("index %d: expected %s, got %s", i, expectedNames[i], m.Name())
		}
	}
}

func TestAptManager(t *testing.T) {
	m := &AptManager{}
	if m.Name() != "apt" {
		t.Errorf("expected apt, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestBrewManager(t *testing.T) {
	m := &BrewManager{}
	if m.Name() != "brew" {
		t.Errorf("expected brew, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestPipManager(t *testing.T) {
	m := &PipManager{}
	if m.Name() != "pip" {
		t.Errorf("expected pip, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestNpmManager(t *testing.T) {
	m := &NpmManager{}
	if m.Name() != "npm" {
		t.Errorf("expected npm, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestCargoManager(t *testing.T) {
	m := &CargoManager{}
	if m.Name() != "cargo" {
		t.Errorf("expected cargo, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestGoManager(t *testing.T) {
	m := &GoManager{}
	if m.Name() != "go" {
		t.Errorf("expected go, got %s", m.Name())
	}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestParseCargoSearchInfo(t *testing.T) {
	input := `serde = "1.0.217" # Serde serialization framework`
	info := parseCargoSearchInfo(input, "serde")
	if info == nil {
		t.Fatal("expected non-nil info")
	}
	if info.Name != "serde" {
		t.Errorf("expected serde, got %s", info.Name)
	}
	if info.Version != "1.0.217" {
		t.Errorf("expected 1.0.217, got %s", info.Version)
	}
}

func TestParseNpmViewJSON(t *testing.T) {
	input := `{"version":"4.18.2","description":"Fast, unopinionated, minimalist web framework","homepage":"https://expressjs.com","license":"MIT"}`
	info, err := parseNpmViewJSON(input, "express")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "express" {
		t.Errorf("expected express, got %s", info.Name)
	}
	if info.Version != "4.18.2" {
		t.Errorf("expected 4.18.2, got %s", info.Version)
	}
	if info.Homepage != "https://expressjs.com" {
		t.Errorf("expected https://expressjs.com, got %s", info.Homepage)
	}
}

func TestParseNpmListJSON(t *testing.T) {
	input := `{"dependencies":{"express":{"version":"4.18.2"},"lodash":{"version":"4.17.21"}}}`
	pkgs, err := parseNpmListJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}
	found := false
	for _, p := range pkgs {
		if p.Name == "express" && p.Version == "4.18.2" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find express 4.18.2")
	}
}
