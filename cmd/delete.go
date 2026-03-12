package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sudabon/dotenv_cryption/internal/config"
)

func newDeleteCmd(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete envcrypt managed resources",
	}

	cmd.AddCommand(newDeleteMasterCmd(deps))

	return cmd
}

func newDeleteMasterCmd(deps Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "master",
		Short: "Delete the configured master secret",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := deps.LoadConfig()
			if err != nil {
				return err
			}

			secretProvider, err := deps.ProviderFactory(cfg)
			if err != nil {
				return err
			}

			if err := secretProvider.DeleteMasterKey(); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted master secret: %s\n", masterSecretID(cfg))
			return nil
		},
	}
}

func masterSecretID(cfg config.Config) string {
	switch cfg.Cloud {
	case "gcp":
		return cfg.GCP.SecretID
	case "aws":
		return cfg.AWS.SecretID
	default:
		return ""
	}
}
