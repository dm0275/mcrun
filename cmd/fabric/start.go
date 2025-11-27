package fabric

import (
	"fmt"
	"strings"

	"github.com/dm0275/mcrun/cmd/common"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewFabricStartCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftFabricConfig()
	fabricCmd := &cobra.Command{
		Use:   "start",
		Short: "Launch a Minecraft Fabric server instance.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Running fabric server: %s", mcConfig.WorldName))

			// Setup directories
			err := minecraft.SetupDirectories(mcConfig)
			utils.CheckErr(err)

			err = minecraft.EnsurePortAvailable(mcConfig.Port, mcConfig.WorldName)
			utils.CheckErr(err)

			if mcConfig.EnableRcon && strings.TrimSpace(mcConfig.RconPort) != "" {
				err = minecraft.EnsurePortAvailable(mcConfig.RconPort, mcConfig.WorldName)
				utils.CheckErr(err)
			}

			// Generate server properties
			if mcConfig.LocalServerConfig {
				err = minecraft.GenerateServerConfig(mcConfig.ServerConfigFile)
				utils.CheckErr(err)
			}

			err = minecraft.SyncMods(mcConfig)
			utils.CheckErr(err)

			err = minecraft.SaveServerMetadata(mcConfig)
			utils.CheckErr(err)

			// Generate Compose file
			err = minecraft.GenerateComposeFile(mcConfig)
			utils.CheckErr(err)

			// Start the server
			err = minecraft.StartServer(mcConfig)
			utils.CheckErr(err)

			fmt.Println("Minecraft server has started ✅")
		},
	}

	// Configure flags
	configureFabricStartFlags(fabricCmd, mcConfig)

	return fabricCmd
}

func configureFabricStartFlags(cmd *cobra.Command, mcconfig *minecraft.MinecraftConfig) {
	common.ConfigureCommonFlags(cmd, mcconfig)
}
