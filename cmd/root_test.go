package cmd

import (
	"bytes"
	"crypto/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/sudabon/dotenv_cryption/internal/config"
	"github.com/sudabon/dotenv_cryption/internal/provider"
)

func TestEncryptDecryptRoundTripWithDefaultPaths(t *testing.T) {
	dir := t.TempDir()
	chdirForTest(t, dir)

	masterKey := mustRandomKey(t)
	writeTestFile(t, filepath.Join(dir, "dotenv.yaml"), `cloud: gcp
gcp:
  project_id: sample-project
  secret_id: sample-secret
`)
	writeTestFile(t, filepath.Join(dir, ".env"), "HELLO=world\n")

	root := newTestRootCmd(t, dir, masterKey)
	root.SetArgs([]string{"encrypt"})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)

	if err := root.Execute(); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	encryptedPath := filepath.Join(dir, ".env.enc")
	assertFileExists(t, encryptedPath)

	root = newTestRootCmd(t, dir, masterKey)
	root.SetArgs([]string{"decrypt"})
	root.SetOut(&stdout)
	root.SetErr(&stdout)

	if err := root.Execute(); err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	assertFileContent(t, filepath.Join(dir, ".env"), "HELLO=world\n")
	if !strings.Contains(stdout.String(), ".env -> .env.enc") {
		t.Fatalf("expected encrypt output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), ".env.enc -> .env") {
		t.Fatalf("expected decrypt output, got %q", stdout.String())
	}
}

func TestEncryptUsesConfiguredPrefix(t *testing.T) {
	dir := t.TempDir()
	chdirForTest(t, dir)

	masterKey := mustRandomKey(t)
	writeTestFile(t, filepath.Join(dir, "dotenv.yaml"), `cloud: aws
aws:
  region: ap-northeast-1
  secret_id: sample-secret
files:
  encrypted_prefix: enc.
`)
	writeTestFile(t, filepath.Join(dir, ".env"), "HELLO=world\n")

	root := newTestRootCmd(t, dir, masterKey)
	root.SetArgs([]string{"encrypt"})

	if err := root.Execute(); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	assertFileExists(t, filepath.Join(dir, "enc..env"))
}

func TestDecryptRejectsInvalidFormat(t *testing.T) {
	dir := t.TempDir()
	chdirForTest(t, dir)

	masterKey := mustRandomKey(t)
	writeTestFile(t, filepath.Join(dir, "dotenv.yaml"), `cloud: gcp
gcp:
  project_id: sample-project
  secret_id: sample-secret
`)
	writeTestFile(t, filepath.Join(dir, ".env.enc"), "plain text")

	root := newTestRootCmd(t, dir, masterKey)
	root.SetArgs([]string{"decrypt"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected decrypt to fail")
	}
	if !strings.Contains(err.Error(), "invalid file format") {
		t.Fatalf("expected invalid format error, got %v", err)
	}
}

func TestCreateMasterCreatesConfiguredSecret(t *testing.T) {
	dir := t.TempDir()
	chdirForTest(t, dir)

	writeTestFile(t, filepath.Join(dir, "dotenv.yaml"), `cloud: gcp
gcp:
  project_id: sample-project
  secret_id: sample-secret
`)

	provider := &trackingProvider{}
	root := newManagedResourceRootCmd(t, dir, provider)
	root.SetArgs([]string{"create", "master"})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)

	if err := root.Execute(); err != nil {
		t.Fatalf("create master failed: %v", err)
	}
	if !provider.createCalled {
		t.Fatal("expected CreateMasterKey to be called")
	}
	if !strings.Contains(stdout.String(), "created master secret: sample-secret") {
		t.Fatalf("expected create output, got %q", stdout.String())
	}
}

func TestDeleteMasterDeletesConfiguredSecret(t *testing.T) {
	dir := t.TempDir()
	chdirForTest(t, dir)

	writeTestFile(t, filepath.Join(dir, "dotenv.yaml"), `cloud: aws
aws:
  region: ap-northeast-1
  secret_id: sample-secret
`)

	provider := &trackingProvider{}
	root := newManagedResourceRootCmd(t, dir, provider)
	root.SetArgs([]string{"delete", "master"})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)

	if err := root.Execute(); err != nil {
		t.Fatalf("delete master failed: %v", err)
	}
	if !provider.deleteCalled {
		t.Fatal("expected DeleteMasterKey to be called")
	}
	if !strings.Contains(stdout.String(), "deleted master secret: sample-secret") {
		t.Fatalf("expected delete output, got %q", stdout.String())
	}
}

func newTestRootCmd(t *testing.T, dir string, masterKey []byte) *cobra.Command {
	t.Helper()

	loadConfig := func() (config.Config, error) {
		return config.LoadFromPath(filepath.Join(dir, "dotenv.yaml"))
	}

	providerFactory := func(config.Config) (provider.SecretProvider, error) {
		return staticProvider{key: masterKey}, nil
	}

	return NewRootCmd(Dependencies{
		LoadConfig:      loadConfig,
		ProviderFactory: providerFactory,
	})
}

func newManagedResourceRootCmd(t *testing.T, dir string, secretProvider provider.SecretProvider) *cobra.Command {
	t.Helper()

	loadConfig := func() (config.Config, error) {
		return config.LoadFromPath(filepath.Join(dir, "dotenv.yaml"))
	}

	providerFactory := func(config.Config) (provider.SecretProvider, error) {
		return secretProvider, nil
	}

	return NewRootCmd(Dependencies{
		LoadConfig:      loadConfig,
		ProviderFactory: providerFactory,
	})
}

type staticProvider struct {
	key []byte
}

func (p staticProvider) GetMasterKey() ([]byte, error) {
	return append([]byte(nil), p.key...), nil
}

func (staticProvider) CreateMasterKey() error {
	return nil
}

func (staticProvider) DeleteMasterKey() error {
	return nil
}

type trackingProvider struct {
	createCalled bool
	deleteCalled bool
}

func (p *trackingProvider) GetMasterKey() ([]byte, error) {
	return nil, nil
}

func (p *trackingProvider) CreateMasterKey() error {
	p.createCalled = true
	return nil
}

func (p *trackingProvider) DeleteMasterKey() error {
	p.deleteCalled = true
	return nil
}

func mustRandomKey(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	return key
}
