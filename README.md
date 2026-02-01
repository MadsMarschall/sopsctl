**Secure Secret Management for SOPS encrypted secrets in FluxCD GitOps repositories**

Sopsctl is a command-line tool that streamlines the management of encrypted Kubernetes secrets in GitOps workflows. It enables developers and DevOps engineers to securely create, edit, and manage secrets in a way that fits naturally into Git-based workflows when using Flux with Sops encrypted secrets stored in git repositories.

## Why Sopsctl?
Adding encrypted secrets to GitOps repositories gives many benefits, including version control, auditability, and collaboration.

https://fluxcd.io/flux/guides/mozilla-sops/

However, managing these secrets can be cumbersome requiring multiple tools and manual steps. When using Flux with SOPS and age encryption, users often need to juggle between `kubectl`, `sops`, and manual key management when editing and creating secrets. Sopsctl simplifies this process by providing a unified CLI that simplifies creating and editing secrets while makeing sure that no unencrypted data is checked into version control.



The tool is invoked using the `sopsctl` command.

## 🌟 Features
* **Create encrypted secrets easily:** Generate SOPS-encrypted Kubernetes secrets from files, literal values, or environment files with a single command.
* **Edit secrets securely:** Edit encrypted secret files with automatic decryption and re-encryption, ensuring sensitive data does not get commited into GitOps repo.
* **Edit individual secret properties in encrypted secret:** Modify specific fields within an encrypted secret without exposing the entire file.
* **Edit encrypted encoded and encrypted secrets:** Seamlessly handle secrets that are both base64-encoded and SOPS-encrypted.

## 📋 Prerequisites

Before using sopsctl, ensure you have:

- **Kubernetes cluster access** with `kubectl` configured. sopsctl retrieves age keys from Kubernetes secrets.
- **A kubernetes cluster with SOPS age keys** set up. Follow the [FluxCD SOPS guide](https://fluxcd.io/flux/guides/mozilla-sops/#encrypting-secrets-using-age) to create and store age keys in your cluster.

## 📦 Installation

### Linux / macOS

Install the latest version using the installation script:

```bash
curl -s https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.sh | sudo bash
```

Or install to a custom directory (no sudo required):

```bash
curl -s https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.sh | bash -s -- ~/.local/bin
```

Install a specific version:

```bash
export SOPSCTL_VERSION=1.0.0
curl -s https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.sh | sudo bash
```

### Windows

Install using PowerShell (run as Administrator or regular user):

```powershell
irm https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.ps1 | iex
```

Install a specific version:

```powershell
$env:Version = "1.0.0"
irm https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.ps1 | iex
```

Install to a custom directory:

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/phantomstacks/sopsctl/main/install/install.ps1))) -BinDir "C:\tools\bin"
```

**Note:** Windows installation requires Windows 10 (version 1803+) or Windows Server 2019+ for the built-in `tar.exe` utility.

### Verify Installation

After installation, verify that sopsctl is working:

```bash
sopsctl --help
```

### Manual Installation

You can also download pre-built binaries from the [releases page](https://github.com/phantomstacks/sopsctl/releases) and manually place them in your PATH.

## 🚀 Quick Start

### 1. Add Age Keys from Your Kubernetes Cluster

First, you need to extract the age keys from your Kubernetes cluster and store them locally:

```bash
# Add keys from your current kubectl context
sopsctl create key --from-cluster

# Or specify a specific cluster context
sopsctl create key --cluster=production

# List stored keys
sopsctl get keys
```
### 2. Create and Manage Encrypted Secrets

```bash
# Create an encrypted secret from literal values
sopsctl create secret my-secret \
  --from-literal=username=admin \
  --from-literal=password=secret123 \
  --cluster=production > secret.yaml

# Create an encrypted secret from an environment file
sopsctl create secret db-credentials \
  --from-env-file=.env \
  --cluster=production > db-secret.yaml

# Edit an encrypted secret file
sopsctl edit secret secrets.yaml --cluster=production

# Decrypt and view a secret (outputs to stdout)
sopsctl get secret secrets.yaml --cluster=production
```

## 📖 Command Reference

### Global Flags

All commands support the following global flag:
- `--cluster, -c`: Specify the Kubernetes cluster context to use for key operations

### Key Management Commands

#### `sopsctl create key`

Add encryption keys from a Kubernetes cluster to local storage. Retrieves age keys from a Kubernetes secret and stores them locally in `~/.sopsctl/` for use with SOPS encryption/decryption.

```bash
sopsctl create key [flags]
```

**Flags:**
- `--from-cluster`: Use the current kubectl context instead of specifying `--cluster`
- `--namespace, -n`: The namespace where the secret is located (default: `flux-system`)
- `--secret, -s`: The name of the secret containing the SOPS key (default: `sops-age`)
- `--key, -k`: The key within the secret that holds the age key (default: `age.agekey`)

**Note:** Either `--from-cluster` or `--cluster` must be specified.

**Examples:**

```bash
# Add keys from current context
sopsctl create key --from-cluster

# Add keys from specific cluster
sopsctl create key --cluster=production

# Add keys from specific cluster and custom secret location
sopsctl create key --cluster=staging --namespace=encryption --secret=my-age-key

# Add keys with custom key name
sopsctl create key --cluster=production --key=private.key
```

#### `sopsctl get keys`

List all age keys stored locally in `~/.sopsctl/`. Shows the cluster name and public key for each stored key.

```bash
sopsctl get keys [flags]
```

**Flags:**
- `--show-sensitive`: Show private keys in the list output (use with caution)

**Examples:**

```bash
# List all stored keys
sopsctl get keys

# List keys including private keys
sopsctl get keys --show-sensitive
```

#### `sopsctl delete key`

Remove age keys from local storage. Can remove keys for a specific cluster or all stored keys.

```bash
sopsctl delete key [cluster-name] [flags]
```

**Flags:**
- `--all`: Remove all SOPS keys from local storage

**Examples:**

```bash
# Remove keys for specific cluster
sopsctl delete key production

# Remove all stored keys
sopsctl delete key --all
```

#### `sopsctl config`

View and manage sopsctl configuration. Controls how and where encryption keys are stored.

```bash
sopsctl config [command]
```

**Subcommands:**
- `view`: Display the current configuration
- `set <key> <value>`: Set a configuration value

**Examples:**

```bash
# View current configuration
sopsctl config view

# Set storage mode to local
sopsctl config set storage-mode local

# Set storage mode to cluster (keys never stored locally)
sopsctl config set storage-mode cluster
```

### Secret Management Commands

#### `sopsctl create secret`

Create an encrypted Kubernetes secret from local files, directories, or literal values. This command mimics `kubectl create secret generic` but automatically encrypts the output using SOPS and age encryption. The secret is output as encrypted YAML that can be committed to version control.

```bash
sopsctl create secret NAME [flags]
```

**Flags:**
- `--from-file strings`: Create secret from files or directories. Can specify `key=filepath` to set custom keys
- `--from-literal stringArray`: Create secret from literal key=value pairs (e.g., `username=admin`)
- `--from-env-file strings`: Create secret from environment files containing `KEY=value` lines
- `--type string`: The type of secret to create (default: `Opaque`)
- `--namespace, -n string`: Namespace for the secret (default: `default`)
- `--append-hash`: Append a hash of the secret data to its name

**Examples:**

```bash
# Create secret from literal values
sopsctl create secret my-secret --from-literal=username=admin --from-literal=password=secret123

# Create secret from files
sopsctl create secret my-secret --from-file=ssh-privatekey=~/.ssh/id_rsa --from-file=ssh-publickey=~/.ssh/id_rsa.pub

# Create secret from a directory (uses filenames as keys)
sopsctl create secret my-secret --from-file=./config/

# Create secret from environment file
sopsctl create secret my-secret --from-env-file=.env

# Create secret with custom namespace and type
sopsctl create secret my-secret --from-literal=token=abc123 --namespace=production --type=kubernetes.io/service-account-token

# Create secret with hash appended to name
sopsctl create secret my-secret --from-literal=data=value --append-hash
```

**Notes:**
- The `--from-env-file` flag cannot be combined with `--from-file` or `--from-literal`
- Output is encrypted SOPS YAML that can be saved to a file: `sopsctl create secret my-secret --from-literal=key=value > secret.yaml`
- Secret data is base64-encoded and then encrypted with SOPS

#### `sopsctl edit secret`

Edit encrypted secret files using your default editor with automatic encryption/decryption. Provides a secure workflow where the file is temporarily decrypted, opened in an editor, then re-encrypted when you save.

```bash
sopsctl edit secret [file] [flags]
```

**Flags:**
- `--decode, -d`: Edit a decoded secret property without manually encrypting the entire file
- `--k, -k string`: Specify the key within the secret to decode and edit (used with `--decode`)
- `--env, -e`: Specify environment variable that holds the decoded value

**Examples:**

```bash
# Edit entire encrypted file
sopsctl edit secret secrets.yaml --cluster=production

# Edit a specific decoded property
sopsctl edit secret secrets.yaml --cluster=production --decode --k=database-password

# Edit single property (auto-detected if only one exists)
sopsctl edit secret secrets.yaml --cluster=production --decode

# Edit with environment variable
sopsctl edit secret secrets.yaml --cluster=production --env
```

**Workflow:**
1. Decrypts the file using the cluster's private age key
2. Opens the decrypted content in your system's default editor
3. After you save and close the editor, re-encrypts the content with the cluster's public key
4. Atomically writes the encrypted content back to the original file

**Editor Selection:**
The command respects the following environment variables (in order of precedence):
1. `SOPSCTL_EDITOR`
2. Default: `nano` (Unix) or `notepad` (Windows)

#### `sopsctl get secret`

Decrypt a SOPS-encrypted file and output the plaintext result to stdout. Useful for viewing encrypted files, piping to other commands, or extracting specific values.

```bash
sopsctl get secret <file> [flags]
```

**Examples:**

```bash
# Decrypt and view file contents
sopsctl get secret secrets.yaml --cluster=production

# Decrypt and pipe to another command
sopsctl get secret secrets.yaml --cluster=production | grep password

# Decrypt and save to a file
sopsctl get secret secrets.yaml --cluster=production > decrypted.yaml

# Use with yq to extract specific values
sopsctl get secret secrets.yaml --cluster=production | yq .data.password
```

**Security Note:** Be careful when decrypting files as the plaintext output may be sensitive. Avoid saving decrypted content to disk unnecessarily.

## ⚙️ Configuration

### Environment Variables

Phantom Flux respects standard environment variables for editor selection:
- `SOPSCTL_EDITOR`: Primary editor preference
- **Default**: `nano` on Unix systems, `notepad` on Windows

### Key Storage

Age keys are stored in `~/.sopsctl/` directory by default. You can change the storage mode using the `sopsctl storage-mode` command.

### SOPS Configuration

Create a `.sops.yaml` file in your project root to configure encryption rules:

```yaml
creation_rules:
  - path_regex: .*secrets.*\.yaml$
    encrypted_regex: ^(data|stringData)$
    age: "age-private-key-here"
  - path_regex: .*config.*\.yaml$
    encrypted_regex: ^(spec\.data)$
    age: "age-private-key-here"
```

## 📝 Common Workflows

### Setting Up a New Environment

1. **Configure age keys in Kubernetes:**
Follow the FluxCD guide to set up SOPS with age encryption in your cluster:
https://fluxcd.io/flux/guides/mozilla-sops/#encrypting-secrets-using-age

2. **Add keys to Phantom Flux:**
   ```bash
   sopsctl create key --from-cluster
   ```
## Troubleshooting

### Key Not Found Error

```
Error: failed to retrieve key from cluster
```

**Solutions:**
- Verify kubectl context: `kubectl config current-context`
- Check secret exists: `kubectl get secret sops-age -n flux-system`
- Ensure you have proper permissions to read secrets

### Editor Not Opening

```
Error: failed to open editor
```

**Solutions:**
- Set `SOPSCTL_EDITOR` environment variable: `export SOPSCTL_EDITOR=nano`
- Ensure your editor is in PATH
- Try with a simple editor like `nano` or `vi`

### SOPS Decryption Failed

```
Error: failed to decrypt file
```

**Solutions:**
- Verify the age key is correctly added: `sopsctl get keys`
- Check SOPS configuration in `.sops.yaml`
- Ensure the file was encrypted with the correct age key

### Permission Denied

```
Error: permission denied when accessing key storage
```

**Solutions:**
- Check permissions on `~/.sopsctl/` directory
- Ensure the directory is owned by your user account
- Recreate the directory: `rm -rf ~/.sopsctl && sopsctl create key --from-cluster`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and enhancement requests.

### Development Setup

```bash
git clone <repository-url>
cd sopsctl
go mod download
go build -o sopsctl .
```

### Running Tests

```bash
go test ./...
```

## 🔗 Related Projects

- [SOPS](https://github.com/mozilla/sops) - Secrets OPerationS
- [Age](https://age-encryption.org/) - Simple, modern encryption
- [Flux](https://fluxcd.io/) - GitOps toolkit for Kubernetes

