package cmd

import (
	"sopsctl/pkg"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource",
	Long: `Edit a resource using the default editor.

Available resources:
  secret    Edit an encrypted secret file`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var editSecretCmd = &cobra.Command{
	Use:   "secret [file]",
	Short: "Edit an encrypted secret file",
	Long: `Edit encrypted secret files using your default editor with automatic encryption/decryption.
Provides a secure workflow where the file is temporarily decrypted, opened in an editor, then re-encrypted when you save.

Workflow:
1. Decrypts the file using the cluster's private age key
2. Opens the decrypted content in your system's default editor
3. After you save and close the editor, re-encrypts the content with the cluster's public key
4. Atomically writes the encrypted content back to the original file

Editor Selection:
The command respects the following environment variables (in order of precedence):
1. SOPSCTL_EDITOR
2. Default: nano (Unix) or notepad (Windows)

GUI editors (VS Code, Sublime Text, Atom, ...) detach from the shell by default, so
sopsctl would re-encrypt the file before you finish editing. Pass the editor's wait
flag explicitly:
  export SOPSCTL_EDITOR="code --wait"     # VS Code / VSCodium / code-insiders
  export SOPSCTL_EDITOR="subl --wait"     # Sublime Text
  export SOPSCTL_EDITOR="atom --wait"     # Atom`,
	Example: `  # Edit entire encrypted file
  sopsctl edit secret secrets.yaml --cluster=production

  # Edit a specific decoded property
  sopsctl edit secret secrets.yaml --cluster=production --decode --k=database-password

  # Edit single property (auto-detected if only one exists)
  sopsctl edit secret secrets.yaml --cluster=production --decode`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.ExecuteCobraCommand(domain.SecretEdit, cmd, args)
	},
}

func init() {
	pkg.InitCobraCommand(domain.SecretEdit, editSecretCmd)

	editCmd.AddCommand(editSecretCmd)
}
