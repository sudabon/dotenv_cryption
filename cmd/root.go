package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sudabon/dotenv_cryption/internal/config"
	"github.com/sudabon/dotenv_cryption/internal/provider"
	"github.com/sudabon/dotenv_cryption/pkg/version"
)

type Dependencies struct {
	LoadConfig      func() (config.Config, error)
	ProviderFactory func(config.Config) (provider.SecretProvider, error)
}

func Execute() error {
	return NewRootCmd(Dependencies{}).Execute()
}

func NewRootCmd(deps Dependencies) *cobra.Command {
	if deps.LoadConfig == nil {
		deps.LoadConfig = config.Load
	}
	if deps.ProviderFactory == nil {
		deps.ProviderFactory = provider.New
	}

	rootCmd := &cobra.Command{
		Use:           "envcrypt",
		Short:         "Encrypt and decrypt .env files with cloud-managed keys",
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newEncryptCmd(deps))
	rootCmd.AddCommand(newDecryptCmd(deps))
	rootCmd.AddCommand(newCreateCmd(deps))
	rootCmd.AddCommand(newDeleteCmd(deps))

	return rootCmd
}
