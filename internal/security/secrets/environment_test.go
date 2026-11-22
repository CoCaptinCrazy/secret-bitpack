package secrets

import (
	"strings"
	"testing"
)

func TestResolveFromEnvironment(t *testing.T) {
	t.Setenv("BITPACK_SECRET_DB_PASSWORD", "value-that-must-not-be-logged")
	values, err := ResolveFromEnvironment([]string{"db-password"})
	if err != nil {
		t.Fatal(err)
	}
	if values["db-password"] != "value-that-must-not-be-logged" {
		t.Fatal("secret value was not resolved")
	}
}

func TestResolveFromEnvironmentReportsMissingNamesOnly(t *testing.T) {
	_, err := ResolveFromEnvironment([]string{"missing-secret"})
	if err == nil || !strings.Contains(err.Error(), "BITPACK_SECRET_MISSING_SECRET") {
		t.Fatalf("expected missing secret name error, got %v", err)
	}
}

func TestResolveRejectsEnvironmentNameCollisions(t *testing.T) {
	t.Setenv("BITPACK_SECRET_DB_PASSWORD", "configured")
	_, err := ResolveFromEnvironment([]string{"db-password", "db_password"})
	if err == nil || !strings.Contains(err.Error(), "same environment variable") {
		t.Fatalf("expected ambiguous secret names to fail, got %v", err)
	}
}
