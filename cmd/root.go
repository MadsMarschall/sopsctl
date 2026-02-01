// Package cmd /*
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sopsctl",
	Short: "Secure configuration management with age encryption and SOPS",
	Long: `sopsctl is a CLI tool for managing encrypted configurations using age keys and SOPS in a Kubernetes cluster.

Features:
  - Age key management with secure storage in ~/.sopsctl
  - SOPS integration for encrypted YAML/JSON files
  - Interactive text editor for editing encrypted secrets
  - Easy key generation and import from Kubernetes clusters

Get started:
  sopsctl create key --from-cluster   # Add age key from current cluster
  sopsctl get keys                    # List stored keys
  sopsctl create secret my-secret     # Create an encrypted secret
  sopsctl --help                      # Show all available commands`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Kubernetes cluster context to use")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(configCmd)
}
