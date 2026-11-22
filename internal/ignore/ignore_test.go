package ignore

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMatches(t *testing.T) {
	if !Matches("*.log", "build/app.log") {
		t.Fatal("expected wildcard to match")
	}
	if Matches("*.log", "build/app.go") {
		t.Fatal("unexpected wildcard match")
	}
	if !Matches("vendor", "vendor") {
		t.Fatal("expected exact match")
	}
}

func TestMatcherSupportsDirectoriesAndNegation(t *testing.T) {
	matcher, err := Parse(strings.NewReader(`node_modules/
*.log
!important.log
/build/
`))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"node_modules", "node_modules/pkg/index.js", "logs/app.log", "build"} {
		if !matcher.Ignored(path, path == "node_modules" || path == "build") {
			t.Errorf("expected %q to be ignored", path)
		}
	}
	for _, path := range []string{"src/main.go", "important.log", "nested/build"} {
		if matcher.Ignored(path, false) {
			t.Errorf("did not expect %q to be ignored", path)
		}
	}
}

func TestLoadUsesSafeDefaultsWithoutIgnoreFile(t *testing.T) {
	matcher, err := Load(filepath.Join(t.TempDir(), ".bpackignore"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{".env", ".bitpack/config.json", "node_modules/pkg/index.js"} {
		if !matcher.Ignored(path, false) {
			t.Errorf("expected safe default to ignore %q", path)
		}
	}
}
