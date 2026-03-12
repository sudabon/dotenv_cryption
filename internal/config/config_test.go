package config

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromPathGCP(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "dotenv.yaml")
	writeTestFile(t, path, `cloud: gcp
gcp:
  project_id: sample-project
  secret_id: sample-secret
`)

	cfg, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath returned error: %v", err)
	}

	if cfg.Cloud != "gcp" {
		t.Fatalf("expected cloud gcp, got %q", cfg.Cloud)
	}
	if cfg.GCP.ProjectID != "sample-project" {
		t.Fatalf("expected project id, got %q", cfg.GCP.ProjectID)
	}
	if cfg.Crypto.Algorithm != AlgorithmAES256GCM {
		t.Fatalf("expected default algorithm %q, got %q", AlgorithmAES256GCM, cfg.Crypto.Algorithm)
	}
}

func TestLoadFromPathMissingFile(t *testing.T) {
	t.Parallel()

	_, err := LoadFromPath(filepath.Join(t.TempDir(), "dotenv.yaml"))
	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestValidateRejectsUnsupportedCloud(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "dotenv.yaml")
	writeTestFile(t, path, `cloud: azure`)

	_, err := LoadFromPath(path)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unsupported cloud provider") {
		t.Fatalf("expected unsupported provider error, got %v", err)
	}
}

func TestValidateRejectsMissingGCPFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "dotenv.yaml")
	writeTestFile(t, path, `cloud: gcp
gcp:
  secret_id: sample-secret
`)

	_, err := LoadFromPath(path)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "project_id") {
		t.Fatalf("expected missing project_id error, got %v", err)
	}
}

func TestValidateRejectsMissingAWSFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "dotenv.yaml")
	writeTestFile(t, path, `cloud: aws
aws:
  secret_id: sample-secret
`)

	_, err := LoadFromPath(path)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "region") {
		t.Fatalf("expected missing region error, got %v", err)
	}
}

func TestPathHelpers(t *testing.T) {
	t.Parallel()

	cfg := Config{}
	if got := cfg.EncryptedPath("/tmp/.env"); got != "/tmp/.env.enc" {
		t.Fatalf("expected suffix path, got %q", got)
	}

	cfg.Files.EncryptedPrefix = "enc."
	if got := cfg.EncryptedPath("/tmp/.env"); got != "/tmp/enc..env" {
		t.Fatalf("expected prefixed path, got %q", got)
	}

	decrypted, err := cfg.DecryptedPath("/tmp/enc..env")
	if err != nil {
		t.Fatalf("DecryptedPath returned error: %v", err)
	}
	if decrypted != "/tmp/.env" {
		t.Fatalf("expected decrypted path /tmp/.env, got %q", decrypted)
	}
}
