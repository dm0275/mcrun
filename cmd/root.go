package cmd

import (
	"github.com/dm0275/mcrun/cmd/forge"
	"github.com/spf13/cobra"
)

type MinecraftRunCLI struct {
	rootCmd *cobra.Command
}

func NewCLI() *MinecraftRunCLI {
	minecraftRunCLI := &MinecraftRunCLI{}
	minecraftRunCLI.initialize()

	// Setup sub-commands
	minecraftRunCLI.rootCmd.AddCommand(forge.NewForgeCmd())

	return minecraftRunCLI
}

func (w *MinecraftRunCLI) initialize() {
	// Initialize the root cmd
	w.rootCmd = &cobra.Command{
		Use:   "mcrun",
		Short: "mcrun us a CLI used to run run dockerized Minecraft servers.",
		Long:  `mcrun us a CLI used to run run dockerized Minecraft servers.`,
	}
}

func (w *MinecraftRunCLI) Execute() {
	cobra.CheckErr(w.rootCmd.Execute())
}
