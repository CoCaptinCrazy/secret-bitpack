package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	blueprint, err := Parse(strings.NewReader(`blueprint "social-api" v1.4
using server
bring "./"
env <|
    PORT = 8080
|>
assemble <|
    npm install
    cargo build --release
|>
verify <|
    npm test
    go test ./...
|>
ignite <|
    npm start
|>
`))
	if err != nil {
		t.Fatal(err)
	}
	if blueprint.Name != "social-api" {
		t.Fatalf("unexpected name: %s", blueprint.Name)
	}
	if blueprint.Version != "1.4" {
		t.Fatalf("unexpected version: %s", blueprint.Version)
	}
	if blueprint.Using != "server" || blueprint.Context != "./" {
		t.Fatalf("unexpected runtime/context: %q, %q", blueprint.Using, blueprint.Context)
	}
	if blueprint.Environment["PORT"] != "8080" {
		t.Fatalf("environment was not parsed: %#v", blueprint.Environment)
	}
	if len(blueprint.Assemble) != 2 || len(blueprint.Verify) != 2 ||
		blueprint.Ignite[0] != "npm start" {
		t.Fatalf("execution commands were not preserved: %#v", blueprint)
	}
}

func TestParseRequiresDeclaration(t *testing.T) {
	if _, err := Parse(strings.NewReader("using server\n")); err == nil {
		t.Fatal("expected missing blueprint declaration to fail")
	}
}

func TestParserPreservesUnknownSectionsAndLanguageCommands(t *testing.T) {
	blueprint, err := Parse(strings.NewReader(`blueprint "mobile-client" v2
	using mobile
	assemble <|
	    flutter pub get
	    flutter build apk
	|>
	verify <|
	    flutter test
	    swift test
	|>
	custom <|
	    generated-code
	|>
	`))
	if err != nil {
		t.Fatal(err)
	}
	if blueprint.Using != "mobile" {
		t.Fatalf("runtime kind was restricted: %q", blueprint.Using)
	}
	if blueprint.Assemble[1] != "flutter build apk" || blueprint.Sections["custom"][0] != "generated-code" {
		t.Fatalf("polyglot commands or custom section were lost: %#v", blueprint)
	}
}

func TestParseExampleBitfile(t *testing.T) {
	path := filepath.Join("..", "..", "..", "examples", "Bitfile")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	blueprint, err := Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	if blueprint.Name != "social-api" || blueprint.Using != "server" ||
		blueprint.Context != "./" || blueprint.Environment["PORT"] != "8080" {
		t.Fatalf("example directives were not parsed: %#v", blueprint)
	}
	if len(blueprint.Assemble) != 3 || len(blueprint.Verify) != 2 ||
		len(blueprint.Ignite) != 1 || len(blueprint.Secrets) != 3 {
		t.Fatalf("example sections were not preserved: %#v", blueprint)
	}
}
