package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func setTempHome(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	oldHome := os.Getenv("HOME")
	oldUserProfile := os.Getenv("USERPROFILE")

	if err := os.Setenv("HOME", tmp); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tmp); err != nil {
		t.Fatalf("set USERPROFILE: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
		_ = os.Setenv("USERPROFILE", oldUserProfile)
	})

	return tmp
}

func TestExpandAndCleanPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}

	got, err := expandAndCleanPath("~")
	if err != nil {
		t.Fatalf("expand '~': %v", err)
	}
	if got != home {
		t.Fatalf("expand '~' = %q, want %q", got, home)
	}

	got, err = expandAndCleanPath("~/notes")
	if err != nil {
		t.Fatalf("expand '~/notes': %v", err)
	}
	want := filepath.Join(home, "notes")
	if got != want {
		t.Fatalf("expand '~/notes' = %q, want %q", got, want)
	}

	got, err = expandAndCleanPath("")
	if err != nil {
		t.Fatalf("expand empty path: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("expand empty path returned non-absolute path: %q", got)
	}
}

func TestTagsHandling(t *testing.T) {
	cfg := &Config{Tags: []string{"a", "b"}}

	replaceTags(cfg, "x, y")
	if strings.Join(cfg.Tags, ",") != "x,y" {
		t.Fatalf("replaceTags = %v", cfg.Tags)
	}

	cfg.Tags = []string{"a", "b"}
	addTags(cfg, "b, c")
	if strings.Join(cfg.Tags, ",") != "a,b,c" {
		t.Fatalf("addTags = %v", cfg.Tags)
	}

	cfg.Tags = []string{"a", "b"}
	mergeTags(cfg, "b, c")
	if strings.Join(cfg.Tags, ",") != "a,b,b,c" {
		t.Fatalf("mergeTags = %v", cfg.Tags)
	}
}

func TestPropertiesHandlingWithNilMap(t *testing.T) {
	cfg := &Config{}

	addProperties(cfg, "a=1")
	if cfg.Properties["a"] != "1" {
		t.Fatalf("addProperties did not set value: %v", cfg.Properties)
	}

	cfg = &Config{}
	mergeProperties(cfg, "b=2")
	if cfg.Properties["b"] != "2" {
		t.Fatalf("mergeProperties did not set value: %v", cfg.Properties)
	}

	cfg = &Config{Properties: map[string]string{"x": "y"}}
	replaceProperties(cfg, "c=3")
	if len(cfg.Properties) != 1 || cfg.Properties["c"] != "3" {
		t.Fatalf("replaceProperties unexpected map: %v", cfg.Properties)
	}
}

func TestCreateFrontmatter(t *testing.T) {
	cfg := &Config{
		Name:       "Test Note",
		Tags:       []string{"alpha", "beta"},
		Properties: map[string]string{"owner": "test"},
	}
	fm := createFrontmatter(cfg)

	if !strings.Contains(fm, "title: Test Note") {
		t.Fatalf("frontmatter missing title: %s", fm)
	}
	if !strings.Contains(fm, "tags: [alpha, beta]") {
		t.Fatalf("frontmatter missing tags: %s", fm)
	}
	if !strings.Contains(fm, "owner: test") {
		t.Fatalf("frontmatter missing property: %s", fm)
	}

	dateRe := regexp.MustCompile(`(?m)^date: \d{4}-\d{2}-\d{2}$`)
	if !dateRe.MatchString(fm) {
		t.Fatalf("frontmatter missing date: %s", fm)
	}
}

func TestLoadConfigAndParseError(t *testing.T) {
	tmpHome := setTempHome(t)
	configDirPath := filepath.Join(tmpHome, configDir)
	if err := os.MkdirAll(configDirPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	goodConfig := []byte("vault_path: ~/vault\noverwrite_mode: fail\n")
	if err := os.WriteFile(filepath.Join(configDirPath, "config"), goodConfig, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, cfgPath, err := loadConfig("config")
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfgPath == nil || *cfgPath == "" {
		t.Fatalf("config path missing")
	}
	if cfg.VaultPath != "~/vault" {
		t.Fatalf("unexpected vault path: %q", cfg.VaultPath)
	}

	badConfig := []byte("vault_path: [bad")
	if err := os.WriteFile(filepath.Join(configDirPath, "bad"), badConfig, 0o644); err != nil {
		t.Fatalf("write bad config: %v", err)
	}

	_, badPath, err := loadConfig("bad")
	if err == nil {
		t.Fatalf("expected error for bad config")
	}
	if badPath == nil || *badPath == "" {
		t.Fatalf("bad config path missing")
	}
}

func TestSaveToObsidianSerialize(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		Name:          "Note",
		VaultPath:     tmpDir,
		OverwriteMode: "fail",
		Properties:    map[string]string{},
	}

	firstPath, err := saveToObsidian("first", cfg)
	if err != nil {
		t.Fatalf("save first: %v", err)
	}
	if _, err := os.Stat(firstPath); err != nil {
		t.Fatalf("first file missing: %v", err)
	}

	cfg.OverwriteMode = "serialize"
	secondPath, err := saveToObsidian("second", cfg)
	if err != nil {
		t.Fatalf("save second: %v", err)
	}
	if secondPath == firstPath {
		t.Fatalf("serialize returned same path: %q", secondPath)
	}
	if _, err := os.Stat(secondPath); err != nil {
		t.Fatalf("second file missing: %v", err)
	}
}
