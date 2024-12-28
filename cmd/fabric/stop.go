package fabric

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewFabricStopCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftFabricConfig()
	fabricCmd := &cobra.Command{
		Use:   "stop",
		Short: "Shut down the  Minecraft Fabric server instance.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Stopping fabric server: %s", mcConfig))

			// Get Compose file
			composeFile, err := minecraft.GetComposeFile(mcConfig)
			utils.CheckErr(err)

			// Start the server
			err = minecraft.StopServer(composeFile)
			utils.CheckErr(err)
		},
	}

	// Configure flags
	configureFabricStopFlags(fabricCmd, mcConfig)

	return fabricCmd
}

func configureFabricStopFlags(cmd *cobra.Command, mcconfig *minecraft.MinecraftConfig) {
	cmd.Flags().StringVarP(&mcconfig.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")
}
