package forge

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/config"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewForgeCmd() *cobra.Command {
	mcConfig := config.NewMinecraftConfig()
	forgeCmd := &cobra.Command{
		Use:   "forge",
		Short: "Run a Forge server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Running forge server: %s", mcConfig))

			// Setup directories
			err := config.SetupDirectories(mcConfig)
			utils.CheckErr(err)

		},
	}

	// Configure flags
	configureFlags(forgeCmd, mcConfig)

	return forgeCmd
}

func configureFlags(cmd *cobra.Command, mcconfig *config.MinecraftConfig) {
	cmd.Flags().StringVarP(&mcconfig.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")

	cmd.Flags().StringVarP(&mcconfig.Version, "version", "", mcconfig.Version, "Minecraft version")
	cmd.Flags().StringVarP(&mcconfig.Port, "port", "", mcconfig.Port, "Server port")
	cmd.Flags().StringVarP(&mcconfig.MinMemory, "min-memory", "", mcconfig.MinMemory, "Minimum memory limit")
	cmd.Flags().StringVarP(&mcconfig.MaxMemory, "max-memory", "", mcconfig.MaxMemory, "Maximum memory limit")
	cmd.Flags().StringVarP(&mcconfig.MaxMemory, "max-memory", "", mcconfig.MaxMemory, "Maximum memory limit")
}
