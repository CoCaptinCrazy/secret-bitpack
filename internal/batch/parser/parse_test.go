package parser

import (
	"strings"
	"testing"
)

func TestParseReleaseBatch(t *testing.T) {
	source := `batch "vice-city-release" v1.0
author:
    Rockstar North
parallel:
    true
failfast:
    true
forge:
    game-client
    game-server
    launcher
verify:
    unit-tests
    integration-tests
publish:
    game-client
    game-server
assemble:
    production.bpack
ignite:
    production
`
	plan, err := Parse(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	if plan.Name != "vice-city-release" || !plan.Parallel || !plan.FailFast {
		t.Fatalf("unexpected batch header: %#v", plan)
	}
	if len(plan.Forge) != 3 || len(plan.Verify) != 2 || len(plan.Publish) != 2 ||
		len(plan.Assemble) != 1 || len(plan.Ignite) != 1 {
		t.Fatalf("batch actions were not parsed: %#v", plan)
	}
}

func TestParseRejectsUnknownSection(t *testing.T) {
	if _, err := Parse(strings.NewReader("batch \"release\" v1\nunknown:\n  x\n")); err == nil {
		t.Fatal("expected unknown batch section to fail")
	}
}
