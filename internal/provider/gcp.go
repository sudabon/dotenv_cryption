package provider

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sudabon/dotenv_cryption/internal/config"
	cryptoutil "github.com/sudabon/dotenv_cryption/internal/crypto"
)

type GCPClient interface {
	AccessSecretVersion(ctx context.Context, name string) ([]byte, error)
	CreateSecret(ctx context.Context, parent, secretID string, data []byte) error
	DeleteSecret(ctx context.Context, name string) error
}

type GCPProvider struct {
	projectID string
	secretID  string
	client    GCPClient
}

func newGCPProvider(cfg config.GCPConfig) (SecretProvider, error) {
	client, err := newGCPClient(context.Background())
	if err != nil {
		return nil, wrapGCPClientError(err)
	}
	return &GCPProvider{
		projectID: cfg.ProjectID,
		secretID:  cfg.SecretID,
		client:    client,
	}, nil
}

func (p *GCPProvider) GetMasterKey() ([]byte, error) {
	key, err := p.client.AccessSecretVersion(context.Background(), p.secretVersionName())
	if err != nil {
		return nil, wrapGCPActionError(p.secretID, "retrieve", err)
	}
	if err := validateMasterKey(key, "GCP Secret Manager"); err != nil {
		return nil, err
	}
	return append([]byte(nil), key...), nil
}

func (p *GCPProvider) CreateMasterKey() error {
	key, err := cryptoutil.GenerateMasterKey()
	if err != nil {
		return err
	}
	if err := p.client.CreateSecret(context.Background(), p.projectName(), p.secretID, key); err != nil {
		return wrapGCPActionError(p.secretID, "create", err)
	}
	return nil
}

func (p *GCPProvider) DeleteMasterKey() error {
	if err := p.client.DeleteSecret(context.Background(), p.secretName()); err != nil {
		return wrapGCPActionError(p.secretID, "delete", err)
	}
	return nil
}

type gcpSDKClient struct {
	client *secretmanager.Client
}

func newGCPClient(ctx context.Context) (GCPClient, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &gcpSDKClient{client: client}, nil
}

func (c *gcpSDKClient) AccessSecretVersion(ctx context.Context, name string) ([]byte, error) {
	resp, err := c.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return resp.Payload.GetData(), nil
}

func (c *gcpSDKClient) CreateSecret(ctx context.Context, parent, secretID string, data []byte) error {
	if _, err := c.client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}); err != nil {
		return err
	}

	name := fmt.Sprintf("%s/secrets/%s", parent, secretID)
	if _, err := c.client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: data,
		},
	}); err != nil {
		cleanupErr := c.DeleteSecret(ctx, name)
		if cleanupErr != nil {
			return fmt.Errorf("add secret version: %w (cleanup failed: %v)", err, cleanupErr)
		}
		return err
	}

	return nil
}

func (c *gcpSDKClient) DeleteSecret(ctx context.Context, name string) error {
	err := c.client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	})
	return err
}

func (p *GCPProvider) projectName() string {
	return fmt.Sprintf("projects/%s", p.projectID)
}

func (p *GCPProvider) secretName() string {
	return fmt.Sprintf("%s/secrets/%s", p.projectName(), p.secretID)
}

func (p *GCPProvider) secretVersionName() string {
	return fmt.Sprintf("%s/versions/latest", p.secretName())
}

func wrapGCPClientError(err error) error {
	message := strings.ToLower(err.Error())

	switch {
	case strings.Contains(message, "credential"),
		strings.Contains(message, "unauth"),
		strings.Contains(message, "permission"),
		strings.Contains(message, "access denied"),
		strings.Contains(message, "application default"):
		return fmt.Errorf("gcp authentication failed: run `gcloud auth application-default login` or set GOOGLE_APPLICATION_CREDENTIALS: %w", err)
	default:
		return fmt.Errorf("failed to initialize gcp secret manager client: %w", err)
	}
}

func wrapGCPActionError(secretID, action string, err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return fmt.Errorf("gcp secret %q not found: %w", secretID, err)
		case codes.AlreadyExists:
			return fmt.Errorf("gcp secret %q already exists: %w", secretID, err)
		case codes.PermissionDenied, codes.Unauthenticated:
			return fmt.Errorf("gcp authentication failed: run `gcloud auth application-default login` or set GOOGLE_APPLICATION_CREDENTIALS: %w", err)
		}
	}

	message := strings.ToLower(err.Error())

	switch {
	case strings.Contains(message, "already exists"):
		return fmt.Errorf("gcp secret %q already exists: %w", secretID, err)
	case strings.Contains(message, "not found"):
		return fmt.Errorf("gcp secret %q not found: %w", secretID, err)
	case strings.Contains(message, "credential"),
		strings.Contains(message, "unauth"),
		strings.Contains(message, "permission"),
		strings.Contains(message, "access denied"),
		strings.Contains(message, "application default"):
		return fmt.Errorf("gcp authentication failed: run `gcloud auth application-default login` or set GOOGLE_APPLICATION_CREDENTIALS: %w", err)
	default:
		switch action {
		case "create":
			return fmt.Errorf("failed to create master key in gcp secret manager: %w", err)
		case "delete":
			return fmt.Errorf("failed to delete master key from gcp secret manager: %w", err)
		default:
			return fmt.Errorf("failed to retrieve master key from gcp secret manager: %w", err)
		}
	}
}
