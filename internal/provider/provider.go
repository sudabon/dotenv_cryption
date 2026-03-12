package provider

import (
	"fmt"

	"github.com/sudabon/dotenv_cryption/internal/config"
)

const masterKeySize = 32

type SecretProvider interface {
	GetMasterKey() ([]byte, error)
	CreateMasterKey() error
	DeleteMasterKey() error
}

type Builders struct {
	GCP func(config.GCPConfig) (SecretProvider, error)
	AWS func(config.AWSConfig) (SecretProvider, error)
}

func New(cfg config.Config) (SecretProvider, error) {
	return NewWithBuilders(cfg, defaultBuilders())
}

func NewWithBuilders(cfg config.Config, builders Builders) (SecretProvider, error) {
	switch cfg.Cloud {
	case "gcp":
		return builders.GCP(cfg.GCP)
	case "aws":
		return builders.AWS(cfg.AWS)
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", cfg.Cloud)
	}
}

func validateMasterKey(key []byte, source string) error {
	if len(key) != masterKeySize {
		return fmt.Errorf("invalid master key from %s: expected %d bytes", source, masterKeySize)
	}
	return nil
}

func defaultBuilders() Builders {
	return Builders{
		GCP: newGCPProvider,
		AWS: newAWSProvider,
	}
}
