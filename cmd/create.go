package cmd

import (
	"sopsctl/pkg"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource (secret or key)",
	Long: `Create a resource such as a secret or key.

Available resources:
  secret    Create an encrypted Kubernetes secret
  key       Add an encryption key from a Kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var createSecretCmd = &cobra.Command{
	Use:   "secret NAME [--type=string] [--from-file=[key=]source] [--from-literal=key1=value1] [--from-env-file=source]",
	Short: "Create an encrypted secret from files, directories, or literal values",
	Long: `Create an encrypted Kubernetes secret from a local file, directory, or literal value.
The secret will be encrypted and output as YAML.

When creating a secret based on a file, the key will default to the basename of the file, 
and the value will default to the file content. If the basename is an invalid key or you 
wish to choose your own, you may specify an alternate key.

When creating a secret based on a directory, each file whose basename is a valid key in 
the directory will be packaged into the secret.`,
	Example: `  # Create a new secret named my-secret with keys for each file in folder bar
  sopsctl create secret my-secret --from-file=path/to/bar

  # Create a new secret named my-secret with specified keys instead of names on disk
  sopsctl create secret my-secret --from-file=ssh-privatekey=path/to/id_rsa --from-file=ssh-publickey=path/to/id_rsa.pub

  # Create a new secret named my-secret with key1=supersecret and key2=topsecret
  sopsctl create secret my-secret --from-literal=key1=supersecret --from-literal=key2=topsecret

  # Create a new secret from env files
  sopsctl create secret my-secret --from-env-file=path/to/foo.env --from-env-file=path/to/bar.env`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.SecretCreate, cmd, args)
	},
}

var createKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Add SOPS keys from a Kubernetes cluster",
	Long: `Add SOPS keys from a Kubernetes cluster to local storage.
Retrieves age keys from a Kubernetes secret and stores them locally for use with SOPS encryption/decryption.
			
You can use the current context with the --from-cluster flag, or specify a context with --cluster.

Either --from-cluster or --cluster <ctx-name> must be specified.`,
	Example: `  # Add keys from current kubectl context
  sopsctl create key --from-cluster

  # Add keys from specific cluster context
  sopsctl create key --cluster=production

  # Add keys from specific cluster and custom secret location
  sopsctl create key --cluster=staging --namespace=encryption --secret=my-age-key`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.KeyAdd, cmd, args)
	},
}

func init() {
	pkg.InitCobraCommand(domain.SecretCreate, createSecretCmd)
	pkg.InitCobraCommand(domain.KeyAdd, createKeyCmd)

	createCmd.AddCommand(createSecretCmd)
	createCmd.AddCommand(createKeyCmd)
}
