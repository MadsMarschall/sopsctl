package set

import (
	"fmt"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

type configSetOptions struct {
	Key   string
	Value string
}

type ConfigSetCmd struct {
	options *configSetOptions
	storage domain.KeyStorage
}

func NewConfigSetCmd(storage domain.KeyStorage) *ConfigSetCmd {
	return &ConfigSetCmd{storage: storage}
}

func (c ConfigSetCmd) InitCmd(cmd *cobra.Command) {
	// No additional flags needed - uses positional args
}

func (c ConfigSetCmd) UseOptions(cmd *cobra.Command, args []string) (domain.CommandExecutor, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("usage: sopsctl config set <key> <value>\n\nAvailable keys:\n  storage-mode    Set the key storage mode (local, cluster)")
	}

	c.options = &configSetOptions{
		Key:   args[0],
		Value: args[1],
	}
	return c, nil
}

func (c ConfigSetCmd) Execute() (string, error) {
	switch c.options.Key {
	case "storage-mode":
		sm := domain.StorageMode(c.options.Value)
		if !sm.IsValid() {
			return "", fmt.Errorf("invalid storage mode: %s (valid options: local, cluster)", c.options.Value)
		}

		currentMode, err := c.storage.GetStorageMode()
		if err != nil {
			return "", err
		}
		if currentMode.ToString() == sm.ToString() {
			return fmt.Sprintf("storage-mode is already set to %s", sm.ToString()), nil
		}

		err = c.storage.SetStorageMode(sm)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("storage-mode set to %s", sm.ToString()), nil

	default:
		return "", fmt.Errorf("unknown configuration key: %s\n\nAvailable keys:\n  storage-mode    Set the key storage mode (local, cluster)", c.options.Key)
	}
}
