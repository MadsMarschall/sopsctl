package view

import (
	"fmt"
	"sopsctl/pkg/domain"

	"github.com/spf13/cobra"
)

type ConfigViewCmd struct {
	storage domain.KeyStorage
}

func NewConfigViewCmd(storage domain.KeyStorage) *ConfigViewCmd {
	return &ConfigViewCmd{storage: storage}
}

func (c ConfigViewCmd) InitCmd(cmd *cobra.Command) {
	// No additional flags needed for view
}

func (c ConfigViewCmd) UseOptions(cmd *cobra.Command, args []string) (domain.CommandExecutor, error) {
	return c, nil
}

func (c ConfigViewCmd) Execute() (string, error) {
	mode, err := c.storage.GetStorageMode()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("storage-mode: %s", mode.ToString()), nil
}
