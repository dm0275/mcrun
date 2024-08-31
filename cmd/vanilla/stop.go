package vanilla

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewStopCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftForgeConfig()
	forgeCmd := &cobra.Command{
		Use:   "stop",
		Short: "Shut down the  Minecraft server instance.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Stopping server: %s", mcConfig))

			// Get Compose file
			composeFile, err := minecraft.GetComposeFile(mcConfig)
			utils.CheckErr(err)

			// Start the server
			err = minecraft.StopServer(composeFile)
			utils.CheckErr(err)
		},
	}

	// Configure flags
	configureForgeStopFlags(forgeCmd, mcConfig)

	return forgeCmd
}

func configureForgeStopFlags(cmd *cobra.Command, mcconfig *minecraft.MinecraftConfig) {
	cmd.Flags().StringVarP(&mcconfig.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")
}
