package servers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	var worldName string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Stop and delete a Minecraft server",
		RunE: func(cmd *cobra.Command, args []string) error {
			worldName = strings.TrimSpace(worldName)
			if worldName == "" {
				return fmt.Errorf("world name is required")
			}

			cfg := minecraft.NewMinecraftConfig()
			cfg.WorldName = worldName
			composeFile, err := minecraft.GetComposeFile(cfg)
			if err == nil {
				if err := minecraft.StopServer(composeFile); err != nil {
					return fmt.Errorf("failed to stop server: %w", err)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to locate compose file: %w", err)
			}

			if err := minecraft.DeleteServerResources(worldName); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted server %s\n", worldName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&worldName, "world-name", "", "", "Name of the server to delete")
	cmd.MarkFlagRequired("world-name")

	return cmd
}
