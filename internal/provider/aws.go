package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awstypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/sudabon/dotenv_cryption/internal/config"
	cryptoutil "github.com/sudabon/dotenv_cryption/internal/crypto"
)

type AWSClient interface {
	GetSecretValue(ctx context.Context, secretID string) ([]byte, error)
	CreateSecret(ctx context.Context, secretID string, data []byte) error
	DeleteSecret(ctx context.Context, secretID string) error
}

type AWSProvider struct {
	region   string
	secretID string
	client   AWSClient
}

func newAWSProvider(cfg config.AWSConfig) (SecretProvider, error) {
	client, err := newAWSClient(context.Background(), cfg.Region)
	if err != nil {
		return nil, wrapAWSClientError(err)
	}
	return &AWSProvider{
		region:   cfg.Region,
		secretID: cfg.SecretID,
		client:   client,
	}, nil
}

func (p *AWSProvider) GetMasterKey() ([]byte, error) {
	key, err := p.client.GetSecretValue(context.Background(), p.secretID)
	if err != nil {
		return nil, wrapAWSError(p.secretID, "retrieve", err)
	}
	if err := validateMasterKey(key, "AWS Secrets Manager"); err != nil {
		return nil, err
	}
	return append([]byte(nil), key...), nil
}

func (p *AWSProvider) CreateMasterKey() error {
	key, err := cryptoutil.GenerateMasterKey()
	if err != nil {
		return err
	}
	if err := p.client.CreateSecret(context.Background(), p.secretID, key); err != nil {
		return wrapAWSError(p.secretID, "create", err)
	}
	return nil
}

func (p *AWSProvider) DeleteMasterKey() error {
	if err := p.client.DeleteSecret(context.Background(), p.secretID); err != nil {
		return wrapAWSError(p.secretID, "delete", err)
	}
	return nil
}

type awsSDKClient struct {
	client *secretsmanager.Client
}

func newAWSClient(ctx context.Context, region string) (AWSClient, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &awsSDKClient{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

func (c *awsSDKClient) GetSecretValue(ctx context.Context, secretID string) ([]byte, error) {
	resp, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}
	if resp.SecretBinary != nil {
		return append([]byte(nil), resp.SecretBinary...), nil
	}
	if resp.SecretString != nil {
		return []byte(*resp.SecretString), nil
	}
	return nil, errors.New("secret has no payload")
}

func (c *awsSDKClient) CreateSecret(ctx context.Context, secretID string, data []byte) error {
	_, err := c.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         &secretID,
		SecretBinary: data,
	})
	return err
}

func (c *awsSDKClient) DeleteSecret(ctx context.Context, secretID string) error {
	forceDeleteWithoutRecovery := true
	_, err := c.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:                   &secretID,
		ForceDeleteWithoutRecovery: &forceDeleteWithoutRecovery,
	})
	return err
}

func wrapAWSClientError(err error) error {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "credential"),
		strings.Contains(message, "accessdenied"),
		strings.Contains(message, "expiredtoken"),
		strings.Contains(message, "unauthorized"):
		return fmt.Errorf("aws authentication failed: configure AWS credentials (for example `aws configure` or AWS_PROFILE): %w", err)
	default:
		return fmt.Errorf("failed to initialize aws secrets manager client: %w", err)
	}
}

func wrapAWSError(secretID, action string, err error) error {
	var notFound *awstypes.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return fmt.Errorf("aws secret %q not found: %w", secretID, err)
	}
	var exists *awstypes.ResourceExistsException
	if errors.As(err, &exists) {
		return fmt.Errorf("aws secret %q already exists: %w", secretID, err)
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "already exists"):
		return fmt.Errorf("aws secret %q already exists: %w", secretID, err)
	case strings.Contains(message, "not found"):
		return fmt.Errorf("aws secret %q not found: %w", secretID, err)
	case strings.Contains(message, "credential"),
		strings.Contains(message, "accessdenied"),
		strings.Contains(message, "expiredtoken"),
		strings.Contains(message, "unauthorized"):
		return fmt.Errorf("aws authentication failed: configure AWS credentials (for example `aws configure` or AWS_PROFILE): %w", err)
	default:
		switch action {
		case "create":
			return fmt.Errorf("failed to create master key in aws secrets manager: %w", err)
		case "delete":
			return fmt.Errorf("failed to delete master key from aws secrets manager: %w", err)
		default:
			return fmt.Errorf("failed to retrieve master key from aws secrets manager: %w", err)
		}
	}
}
