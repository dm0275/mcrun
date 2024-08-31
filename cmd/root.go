package cmd

import (
	"github.com/dm0275/mcrun/cmd/forge"
	"github.com/dm0275/mcrun/cmd/setup"
	"github.com/spf13/cobra"
)

type MinecraftRunCLI struct {
	rootCmd *cobra.Command
}

func NewCLI() *MinecraftRunCLI {
	minecraftRunCLI := &MinecraftRunCLI{}
	minecraftRunCLI.initialize()

	// Setup sub-commands
	minecraftRunCLI.rootCmd.AddCommand(setup.NewSetupCmd())
	minecraftRunCLI.rootCmd.AddCommand(forge.NewForgeCmd())
	minecraftRunCLI.rootCmd.AddCommand(NewVersionCmd())

	minecraftRunCLI.rootCmd.CompletionOptions.DisableDefaultCmd = true

	return minecraftRunCLI
}

func (m *MinecraftRunCLI) initialize() {
	// Initialize the root cmd
	m.rootCmd = &cobra.Command{
		Use:   "mcrun",
		Short: "mcrun is a CLI utility for creating Minecraft servers.",
		Long:  `mcrun is a command-line interface (CLI) utility for creating Minecraft servers.`,
	}
}

func (m *MinecraftRunCLI) Execute() {
	cobra.CheckErr(m.rootCmd.Execute())
}
