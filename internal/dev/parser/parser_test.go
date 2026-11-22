package parser

import (
	"strings"
	"testing"
)

func TestParseDevelopmentFile(t *testing.T) {
	source := `development "social-api" v1.0
using node20
workspace "./"
service postgres = postgres -D .bitpack/dev/postgres
service redis = redis-server
port 3000
watch src/
env :|
    APP_ENV = development
    DEBUG = true
|:
startup :|
    npm install
|:
shell :|
    npm run dev
|:
secrets :|
    db-password
|:
`
	environment, err := Parse(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	if environment.Name != "social-api" || environment.Runtime != "node20" {
		t.Fatalf("unexpected environment: %#v", environment)
	}
	if len(environment.Services) != 2 || environment.Ports[0] != 3000 ||
		environment.Services[0].Command != "postgres -D .bitpack/dev/postgres" {
		t.Fatalf("unexpected services or ports: %#v", environment)
	}
	if environment.Variables["DEBUG"] != "true" || environment.Startup[0] != "npm install" ||
		environment.Shell[0] != "npm run dev" {
		t.Fatalf("development sections were not parsed: %#v", environment)
	}
	if len(environment.Secrets) != 1 || environment.Secrets[0] != "db-password" {
		t.Fatalf("secret names were not parsed: %#v", environment.Secrets)
	}
}

func TestParseRejectsInvalidPortAndUnclosedSection(t *testing.T) {
	if _, err := Parse(strings.NewReader("development \"api\" v1\nport 70000\n")); err == nil {
		t.Fatal("expected invalid port to fail")
	}
	if _, err := Parse(strings.NewReader("development \"api\" v1\nenv :|\nA=1\n")); err == nil {
		t.Fatal("expected unclosed section to fail")
	}
}
