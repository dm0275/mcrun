package forge

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Config struct {
}

func NewForgeCmd() *cobra.Command {
	config := Config{}
	forgeCmd := &cobra.Command{
		Use:   "forge",
		Short: "Run a Forge server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running a forge server")
		},
	}

	// Configure flags
	configureFlags(forgeCmd, config)

	return forgeCmd
}

func configureFlags(cmd *cobra.Command, config Config) {

}
