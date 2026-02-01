package cmd

import (
	"sopsctl/pkg"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify sopsctl configuration",
	Long: `Modify sopsctl configuration files and settings.

Available subcommands:
  view      Display the current configuration
  set       Set a configuration value`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to 'view' behavior when no subcommand is specified
		pkg.ExecuteCobraCommand(domain.ConfigView, cmd, args)
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display the current configuration",
	Long:  `Display the current sopsctl configuration including storage mode and other settings.`,
	Example: `  # View current configuration
  sopsctl config view`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.ConfigView, cmd, args)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value in sopsctl.

Available configuration keys:
  storage-mode    Set the key storage mode (local, cluster)`,
	Example: `  # Set storage mode to local
  sopsctl config set storage-mode local

  # Set storage mode to cluster (keys never stored locally)
  sopsctl config set storage-mode cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.ConfigSet, cmd, args)
	},
}

func init() {
	pkg.InitCobraCommand(domain.ConfigView, configViewCmd)
	pkg.InitCobraCommand(domain.ConfigSet, configSetCmd)

	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configSetCmd)
}
