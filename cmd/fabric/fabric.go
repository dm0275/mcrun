package fabric

import (
	"github.com/spf13/cobra"
)

func NewFabricCmd() *cobra.Command {
	fabricCmd := &cobra.Command{
		Use:   "fabric",
		Short: "Configures a Minecraft Fabric server instance.",
	}

	// Add sub-commands
	fabricCmd.AddCommand(NewFabricStartCmd())
	fabricCmd.AddCommand(NewFabricStopCmd())

	return fabricCmd
}
