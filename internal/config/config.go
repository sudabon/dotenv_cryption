package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	AlgorithmAES256GCM = "aes-256-gcm"
	defaultSuffix      = ".enc"
)

type Config struct {
	Cloud  string       `mapstructure:"cloud"`
	AWS    AWSConfig    `mapstructure:"aws"`
	GCP    GCPConfig    `mapstructure:"gcp"`
	Crypto CryptoConfig `mapstructure:"crypto"`
	Files  FilesConfig  `mapstructure:"files"`
}

type AWSConfig struct {
	Region   string `mapstructure:"region"`
	SecretID string `mapstructure:"secret_id"`
}

type GCPConfig struct {
	ProjectID string `mapstructure:"project_id"`
	SecretID  string `mapstructure:"secret_id"`
}

type CryptoConfig struct {
	Algorithm string `mapstructure:"algorithm"`
}

type FilesConfig struct {
	EncryptedPrefix string `mapstructure:"encrypted_prefix"`
}

var ErrConfigNotFound = errors.New("dotenv.yaml not found")

func Load() (Config, error) {
	return LoadFromPath("dotenv.yaml")
}

func LoadFromPath(path string) (Config, error) {
	var cfg Config

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, notFoundError(path)
		}
		return cfg, fmt.Errorf("failed to read %s: %w", filepath.Base(path), err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetDefault("crypto.algorithm", AlgorithmAES256GCM)
	v.SetDefault("files.encrypted_prefix", "")

	if err := v.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to read %s: %w", filepath.Base(path), err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse %s: %w", filepath.Base(path), err)
	}

	cfg.Cloud = strings.ToLower(strings.TrimSpace(cfg.Cloud))
	cfg.Crypto.Algorithm = strings.ToLower(strings.TrimSpace(cfg.Crypto.Algorithm))
	cfg.Files.EncryptedPrefix = strings.TrimSpace(cfg.Files.EncryptedPrefix)

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	switch c.Cloud {
	case "gcp":
		if missing := missingFields(
			requiredField("project_id", c.GCP.ProjectID),
			requiredField("secret_id", c.GCP.SecretID),
		); len(missing) > 0 {
			return fmt.Errorf("missing required gcp fields: %s", strings.Join(missing, ", "))
		}
	case "aws":
		if missing := missingFields(
			requiredField("region", c.AWS.Region),
			requiredField("secret_id", c.AWS.SecretID),
		); len(missing) > 0 {
			return fmt.Errorf("missing required aws fields: %s", strings.Join(missing, ", "))
		}
	default:
		return fmt.Errorf("unsupported cloud provider: %s", c.Cloud)
	}

	if c.Crypto.Algorithm == "" {
		c.Crypto.Algorithm = AlgorithmAES256GCM
	}
	if c.Crypto.Algorithm != AlgorithmAES256GCM {
		return fmt.Errorf("unsupported crypto algorithm: %s", c.Crypto.Algorithm)
	}

	return nil
}

func (c Config) EncryptedPath(inputPath string) string {
	if c.Files.EncryptedPrefix != "" {
		dir := filepath.Dir(inputPath)
		base := filepath.Base(inputPath)
		return filepath.Join(dir, c.Files.EncryptedPrefix+base)
	}
	return inputPath + defaultSuffix
}

func (c Config) DecryptedPath(inputPath string) (string, error) {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)

	if c.Files.EncryptedPrefix != "" && strings.HasPrefix(base, c.Files.EncryptedPrefix) {
		return filepath.Join(dir, strings.TrimPrefix(base, c.Files.EncryptedPrefix)), nil
	}
	if strings.HasSuffix(base, defaultSuffix) {
		return filepath.Join(dir, strings.TrimSuffix(base, defaultSuffix)), nil
	}

	return "", fmt.Errorf("cannot derive output path from %q", inputPath)
}

func notFoundError(path string) error {
	if filepath.Base(path) == "dotenv.yaml" {
		return ErrConfigNotFound
	}
	return fmt.Errorf("%s not found", filepath.Base(path))
}

func requiredField(name string, value string) string {
	if strings.TrimSpace(value) == "" {
		return name
	}
	return ""
}

func missingFields(fields ...string) []string {
	var missing []string
	for _, field := range fields {
		if field != "" {
			missing = append(missing, field)
		}
	}
	return missing
}
