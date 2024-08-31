package forge

import (
	"github.com/spf13/cobra"
)

func NewForgeCmd() *cobra.Command {
	forgeCmd := &cobra.Command{
		Use:   "forge",
		Short: "Configures a Minecraft Forge server instance.",
	}

	// Add sub-commands
	forgeCmd.AddCommand(NewForgeStartCmd())
	forgeCmd.AddCommand(NewForgeStopCmd())

	return forgeCmd
}
