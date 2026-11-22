package cli

import (
	"reflect"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	got, err := splitCommand(`verify "folder/my blueprint.bit"`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"verify", "folder/my blueprint.bit"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCommand() = %#v, want %#v", got, want)
	}
}

func TestSplitCommandPreservesWindowsPath(t *testing.T) {
	got, err := splitCommand(`verify C:\blueprints\Bitfile`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"verify", `C:\blueprints\Bitfile`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCommand() = %#v, want %#v", got, want)
	}
}

func TestSplitCommandRejectsUnclosedQuote(t *testing.T) {
	if _, err := splitCommand(`verify "Bitfile`); err == nil {
		t.Fatal("expected an unclosed quote to fail")
	}
}
