package servers

import (
	"fmt"

	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List provisioned Minecraft servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			servers, err := minecraft.ListServers()
			if err != nil {
				return err
			}

			if len(servers) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No servers found.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%-18s %-8s %-10s %-12s %-6s %-8s %s\n", "WORLD", "STATUS", "TYPE", "VERSION", "MODS", "COMPOSE", "PATH")
			for _, s := range servers {
				compose := "no"
				if s.HasCompose {
					compose = "yes"
				}
				serverType := "-"
				version := "-"
				modCount := 0
				if s.Metadata != nil {
					if s.Metadata.Type != "" {
						serverType = s.Metadata.Type
					}
					if s.Metadata.Version != "" {
						version = s.Metadata.Version
					}
					modCount = len(s.Metadata.Mods)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-18s %-8s %-10s %-12s %-6d %-8s %s\n", s.WorldName, s.Status, serverType, version, modCount, compose, s.Path)
			}

			return nil
		},
	}

	return cmd
}
