package cmd

import (
	"sopsctl/pkg"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or more resources",
	Long: `Display one or more resources.

Available resources:
  keys      List all stored encryption keys
  secret    Decrypt and display a secret file`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var getKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List all stored SOPS keys",
	Long:  `List all SOPS age keys stored locally in ~/.sopsctl/. Shows the cluster name and public key for each stored key.`,
	Example: `  # List all stored keys
  sopsctl get keys

  # List keys including private keys
  sopsctl get keys --show-sensitive`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.KeyList, cmd, args)
	},
}

var getSecretCmd = &cobra.Command{
	Use:   "secret <file>",
	Short: "Decrypt and display a secret file",
	Long: `Decrypt a SOPS-encrypted file and output the plaintext result to stdout.
Useful for viewing encrypted files, piping to other commands, or extracting specific values.`,
	Example: `  # Decrypt and view file contents
  sopsctl get secret secrets.yaml --cluster=production

  # Decrypt and pipe to another command
  sopsctl get secret secrets.yaml --cluster=production | grep password

  # Use with yq to extract specific values
  sopsctl get secret secrets.yaml --cluster=production | yq .data.password`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.SecretDecrypt, cmd, args)
	},
}

func init() {
	pkg.InitCobraCommand(domain.KeyList, getKeysCmd)
	pkg.InitCobraCommand(domain.SecretDecrypt, getSecretCmd)

	getCmd.AddCommand(getKeysCmd)
	getCmd.AddCommand(getSecretCmd)
}
