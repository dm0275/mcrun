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

			fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-8s %s\n", "WORLD", "COMPOSE", "PATH")
			for _, s := range servers {
				compose := "no"
				if s.HasCompose {
					compose = "yes"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-8s %s\n", s.WorldName, compose, s.Path)
			}

			return nil
		},
	}

	return cmd
}
