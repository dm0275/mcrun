package setup

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/utils"
	"github.com/spf13/cobra"
)

func NewSetupCmd() *cobra.Command {
	mcConfig := minecraft.NewMinecraftConfig()
	forgeCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup Minecraft server directory structure Forge server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Running forge server: %s", mcConfig))
			err := minecraft.SetupDirectories(mcConfig)
			utils.CheckErr(err)
		},
	}

	// Configure flags
	configureFlags(forgeCmd, mcConfig)

	return forgeCmd
}

func configureFlags(cmd *cobra.Command, config *minecraft.MinecraftConfig) {
	cmd.Flags().StringVarP(&config.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")
}
