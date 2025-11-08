package servers

import "github.com/spf13/cobra"

func NewServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "Manage Minecraft servers",
	}

	cmd.AddCommand(NewListCmd())

	return cmd
}
