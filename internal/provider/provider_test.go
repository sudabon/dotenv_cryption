package provider

import (
	"context"
	"errors"
	"strings"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sudabon/dotenv_cryption/internal/config"
)

func TestGCPProviderGetMasterKey(t *testing.T) {
	t.Parallel()

	p := &GCPProvider{
		projectID: "sample-project",
		secretID:  "sample-secret",
		client:    &mockGCPClient{accessData: bytesOfLength(masterKeySize)},
	}

	key, err := p.GetMasterKey()
	if err != nil {
		t.Fatalf("GetMasterKey returned error: %v", err)
	}
	if len(key) != masterKeySize {
		t.Fatalf("expected %d byte key, got %d", masterKeySize, len(key))
	}
}

func TestGCPProviderReturnsAuthGuidance(t *testing.T) {
	t.Parallel()

	p := &GCPProvider{
		projectID: "sample-project",
		secretID:  "sample-secret",
		client:    &mockGCPClient{accessErr: errors.New("application default credentials are missing")},
	}

	_, err := p.GetMasterKey()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "gcloud auth application-default login") {
		t.Fatalf("expected auth guidance, got %v", err)
	}
}

func TestGCPProviderCreateMasterKey(t *testing.T) {
	t.Parallel()

	client := &mockGCPClient{}
	p := &GCPProvider{
		projectID: "sample-project",
		secretID:  "sample-secret",
		client:    client,
	}

	if err := p.CreateMasterKey(); err != nil {
		t.Fatalf("CreateMasterKey returned error: %v", err)
	}
	if client.createdParent != "projects/sample-project" {
		t.Fatalf("expected parent path to be recorded, got %q", client.createdParent)
	}
	if client.createdSecretID != "sample-secret" {
		t.Fatalf("expected secret id to be recorded, got %q", client.createdSecretID)
	}
	if len(client.createdData) != masterKeySize {
		t.Fatalf("expected %d byte key, got %d", masterKeySize, len(client.createdData))
	}
}

func TestGCPProviderCreateMasterKeyReturnsAlreadyExists(t *testing.T) {
	t.Parallel()

	p := &GCPProvider{
		projectID: "sample-project",
		secretID:  "sample-secret",
		client:    &mockGCPClient{createErr: status.Error(codes.AlreadyExists, "exists")},
	}

	err := p.CreateMasterKey()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already exists error, got %v", err)
	}
}

func TestGCPProviderDeleteMasterKey(t *testing.T) {
	t.Parallel()

	client := &mockGCPClient{}
	p := &GCPProvider{
		projectID: "sample-project",
		secretID:  "sample-secret",
		client:    client,
	}

	if err := p.DeleteMasterKey(); err != nil {
		t.Fatalf("DeleteMasterKey returned error: %v", err)
	}
	if client.deletedName != "projects/sample-project/secrets/sample-secret" {
		t.Fatalf("expected delete path to be recorded, got %q", client.deletedName)
	}
}

func TestAWSProviderGetMasterKey(t *testing.T) {
	t.Parallel()

	p := &AWSProvider{
		region:   "ap-northeast-1",
		secretID: "sample-secret",
		client:   &mockAWSClient{accessData: bytesOfLength(masterKeySize)},
	}

	key, err := p.GetMasterKey()
	if err != nil {
		t.Fatalf("GetMasterKey returned error: %v", err)
	}
	if len(key) != masterKeySize {
		t.Fatalf("expected %d byte key, got %d", masterKeySize, len(key))
	}
}

func TestAWSProviderReturnsNotFoundError(t *testing.T) {
	t.Parallel()

	p := &AWSProvider{
		region:   "ap-northeast-1",
		secretID: "sample-secret",
		client:   &mockAWSClient{accessErr: &awstypes.ResourceNotFoundException{Message: awsString("missing")}},
	}

	_, err := p.GetMasterKey()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestAWSProviderCreateMasterKey(t *testing.T) {
	t.Parallel()

	client := &mockAWSClient{}
	p := &AWSProvider{
		region:   "ap-northeast-1",
		secretID: "sample-secret",
		client:   client,
	}

	if err := p.CreateMasterKey(); err != nil {
		t.Fatalf("CreateMasterKey returned error: %v", err)
	}
	if client.createdSecretID != "sample-secret" {
		t.Fatalf("expected secret id to be recorded, got %q", client.createdSecretID)
	}
	if len(client.createdData) != masterKeySize {
		t.Fatalf("expected %d byte key, got %d", masterKeySize, len(client.createdData))
	}
}

func TestAWSProviderCreateMasterKeyReturnsAlreadyExists(t *testing.T) {
	t.Parallel()

	p := &AWSProvider{
		region:   "ap-northeast-1",
		secretID: "sample-secret",
		client:   &mockAWSClient{createErr: &awstypes.ResourceExistsException{Message: awsString("exists")}},
	}

	err := p.CreateMasterKey()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already exists error, got %v", err)
	}
}

func TestAWSProviderDeleteMasterKey(t *testing.T) {
	t.Parallel()

	client := &mockAWSClient{}
	p := &AWSProvider{
		region:   "ap-northeast-1",
		secretID: "sample-secret",
		client:   client,
	}

	if err := p.DeleteMasterKey(); err != nil {
		t.Fatalf("DeleteMasterKey returned error: %v", err)
	}
	if client.deletedSecretID != "sample-secret" {
		t.Fatalf("expected secret id to be recorded, got %q", client.deletedSecretID)
	}
}

func TestNewWithBuildersSelectsConfiguredProvider(t *testing.T) {
	t.Parallel()

	gcpCalled := false
	awsCalled := false

	_, err := NewWithBuilders(config.Config{
		Cloud: "gcp",
		GCP: config.GCPConfig{
			ProjectID: "sample-project",
			SecretID:  "sample-secret",
		},
	}, Builders{
		GCP: func(config.GCPConfig) (SecretProvider, error) {
			gcpCalled = true
			return staticProvider{}, nil
		},
		AWS: func(config.AWSConfig) (SecretProvider, error) {
			awsCalled = true
			return staticProvider{}, nil
		},
	})
	if err != nil {
		t.Fatalf("NewWithBuilders returned error: %v", err)
	}
	if !gcpCalled {
		t.Fatal("expected GCP builder to be called")
	}
	if awsCalled {
		t.Fatal("did not expect AWS builder to be called")
	}
}

type mockGCPClient struct {
	accessData      []byte
	accessErr       error
	createErr       error
	deleteErr       error
	createdParent   string
	createdSecretID string
	createdData     []byte
	deletedName     string
}

func (c *mockGCPClient) AccessSecretVersion(context.Context, string) ([]byte, error) {
	if c.accessErr != nil {
		return nil, c.accessErr
	}
	return c.accessData, nil
}

func (c *mockGCPClient) CreateSecret(_ context.Context, parent, secretID string, data []byte) error {
	c.createdParent = parent
	c.createdSecretID = secretID
	c.createdData = append([]byte(nil), data...)
	if c.createErr != nil {
		return c.createErr
	}
	return nil
}

func (c *mockGCPClient) DeleteSecret(_ context.Context, name string) error {
	c.deletedName = name
	if c.deleteErr != nil {
		return c.deleteErr
	}
	return nil
}

type mockAWSClient struct {
	accessData      []byte
	accessErr       error
	createErr       error
	deleteErr       error
	createdSecretID string
	createdData     []byte
	deletedSecretID string
}

func (c *mockAWSClient) GetSecretValue(context.Context, string) ([]byte, error) {
	if c.accessErr != nil {
		return nil, c.accessErr
	}
	return c.accessData, nil
}

func (c *mockAWSClient) CreateSecret(_ context.Context, secretID string, data []byte) error {
	c.createdSecretID = secretID
	c.createdData = append([]byte(nil), data...)
	if c.createErr != nil {
		return c.createErr
	}
	return nil
}

func (c *mockAWSClient) DeleteSecret(_ context.Context, secretID string) error {
	c.deletedSecretID = secretID
	if c.deleteErr != nil {
		return c.deleteErr
	}
	return nil
}

type staticProvider struct{}

func (staticProvider) GetMasterKey() ([]byte, error) {
	return bytesOfLength(masterKeySize), nil
}

func (staticProvider) CreateMasterKey() error {
	return nil
}

func (staticProvider) DeleteMasterKey() error {
	return nil
}

func bytesOfLength(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i + 1)
	}
	return data
}

func awsString(value string) *string {
	return &value
}
