package vanilla

import (
	"fmt"
	"github.com/dm0275/mcrun/cmd/common"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftConfig()
	forgeCmd := &cobra.Command{
		Use:   "start",
		Short: "Launch a Minecraft server instance.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Starting server: %s", mcConfig.WorldName))

			// Setup directories
			err := minecraft.SetupDirectories(mcConfig)
			utils.CheckErr(err)

			// Generate Compose file
			err = minecraft.GenerateComposeFile(mcConfig)
			utils.CheckErr(err)

			// Start the server
			err = minecraft.StartServer(mcConfig)
			utils.CheckErr(err)
		},
	}

	// Configure flags
	configureForgeStartFlags(forgeCmd, mcConfig)

	return forgeCmd
}

func configureForgeStartFlags(cmd *cobra.Command, mcconfig *minecraft.MinecraftConfig) {
	common.ConfigureCommonFlags(cmd, mcconfig)
}
