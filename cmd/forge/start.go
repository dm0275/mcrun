package forge

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewForgeStartCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftForgeConfig()
	forgeCmd := &cobra.Command{
		Use:   "start",
		Short: "Launch a Minecraft Forge server instance.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Running forge server: %s", mcConfig))

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
	cmd.Flags().StringVarP(&mcconfig.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")

	cmd.Flags().StringVarP(&mcconfig.Version, "version", "", mcconfig.Version, "Minecraft version")
	cmd.Flags().StringVarP(&mcconfig.Port, "port", "", mcconfig.Port, "Server port")
	cmd.Flags().StringVarP(&mcconfig.MinMemory, "min-memory", "", mcconfig.MinMemory, "Minimum memory limit")
	cmd.Flags().StringVarP(&mcconfig.MaxMemory, "max-memory", "", mcconfig.MaxMemory, "Maximum memory limit")
	cmd.Flags().StringVarP(&mcconfig.Seed, "seed", "", mcconfig.Seed, "Minecraft Seed")
	cmd.Flags().StringVarP(&mcconfig.GameMode, "gamemode", "", "0", "Gamemode: Survival mode is gametype=0, Creative is gametype=1, Adventure is gametype=2, and Spectator is gametype=3")
}
