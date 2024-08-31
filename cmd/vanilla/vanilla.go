package vanilla

import "github.com/spf13/cobra"

func NewVanillaCmd() *cobra.Command {
	forgeCmd := &cobra.Command{
		Use:   "vanilla",
		Short: "Configures a Minecraft server (vanilla) instance.",
	}

	// Add sub-commands
	forgeCmd.AddCommand(NewStartCmd())
	forgeCmd.AddCommand(NewStopCmd())

	return forgeCmd
}
