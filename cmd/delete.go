package cmd

import (
	"sopsctl/pkg"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long: `Delete resources by name.

Available resources:
  key       Remove an encryption key from local storage`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var deleteKeyCmd = &cobra.Command{
	Use:   "key [cluster-name]",
	Short: "Remove SOPS keys from local storage",
	Long:  `Remove age keys from local storage. Can remove keys for a specific cluster or all stored keys.`,
	Example: `  # Remove keys for specific cluster
  sopsctl delete key production

  # Remove all stored keys
  sopsctl delete key --all
  sopsctl delete key -a`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.KeyRemove, cmd, args)
	},
}

func init() {
	pkg.InitCobraCommand(domain.KeyRemove, deleteKeyCmd)

	deleteCmd.AddCommand(deleteKeyCmd)
}
