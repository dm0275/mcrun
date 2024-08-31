package cmd

import (
	"fmt"
	"github.com/dm0275/mcrun/pkg/version"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	forgeCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the current version of the mcrun CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	}

	return forgeCmd
}
