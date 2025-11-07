package common

import (
	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/spf13/cobra"
)

func ConfigureCommonFlags(cmd *cobra.Command, mcconfig *minecraft.MinecraftConfig) {
	cmd.Flags().StringVarP(&mcconfig.WorldName, "world-name", "", "", "Name for the Minecraft server")
	cmd.MarkFlagRequired("world-name")

	cmd.Flags().StringVarP(&mcconfig.Version, "version", "", mcconfig.Version, "Minecraft version")
	cmd.Flags().StringVarP(&mcconfig.Port, "port", "", mcconfig.Port, "Server port")
	cmd.Flags().StringVarP(&mcconfig.MinMemory, "min-memory", "", mcconfig.MinMemory, "Minimum memory limit")
	cmd.Flags().StringVarP(&mcconfig.MaxMemory, "max-memory", "", mcconfig.MaxMemory, "Maximum memory limit")
	cmd.Flags().StringVarP(&mcconfig.Seed, "seed", "", mcconfig.Seed, "Minecraft Seed")
	cmd.Flags().BoolVarP(&mcconfig.EnableCmdBlock, "enable-cmd-block", "", mcconfig.EnableCmdBlock, "Enable Command block")
	cmd.Flags().StringVarP(&mcconfig.GameMode, "gamemode", "", mcconfig.GameMode, "Gamemode: Survival mode is gametype=0, Creative is gametype=1, Adventure is gametype=2, and Spectator is gametype=3")
	cmd.Flags().StringSliceVarP(&mcconfig.MountDirs, "mound-dirs", "", []string{}, "Specifies the local directories that should be mounted")
	cmd.Flags().StringVar(&mcconfig.CurseForgeAPIKey, "curseforge-api-key", mcconfig.CurseForgeAPIKey, "API key for CurseForge downloads (defaults to CURSEFORGE_API_KEY env var)")
	cmd.Flags().Var(newCurseForgeModFlag(&mcconfig.Mods), "curseforge-mod", "Fetch a mod from CurseForge using projectID:fileID (repeatable)")
}
