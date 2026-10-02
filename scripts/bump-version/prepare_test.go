package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareOnlyUpdatesSourceWithoutPublishing(t *testing.T) {
	r := releaseRepo(t)
	before, err := r.git("rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"bump-version", "-root", r.root, "-prepare-only", "-date", "2026-10-02", "-nix-command", "must-not-run", "9.9.9"}
	if err := Main(t.Context(), args, io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	if err := Check(r.root, targets, io.Discard); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(r.root, changelogFile))
	if err != nil || !bytes.Contains(content, []byte("/v9.9.9) - 2026-10-02")) {
		t.Fatalf("fixed changelog date missing: %v", err)
	}
	if after, err := r.git("rev-parse", "HEAD"); err != nil || after != before {
		t.Fatalf("preparation committed: %q, %v", after, err)
	}
	if tags, err := r.git("tag", "--list"); err != nil || tags != "" {
		t.Fatalf("preparation tagged: %q, %v", tags, err)
	}
	if refs, err := r.git("ls-remote", "origin"); err != nil || refs != "" {
		t.Fatalf("preparation pushed: %q, %v", refs, err)
	}
}

func TestPrepareOnlyRejectsPublicationAndMissingDate(t *testing.T) {
	for _, extra := range [][]string{
		{},
		{"-date", "2026-10-02", "-commit"},
		{"-date", "2026-10-02", "-push"},
		{"-date", "2026-02-30"},
	} {
		args := append([]string{"bump-version", "-prepare-only"}, extra...)
		args = append(args, "9.9.9")
		if err := Main(t.Context(), args, io.Discard, io.Discard); err == nil {
			t.Fatalf("accepted %v", args)
		}
	}
}
