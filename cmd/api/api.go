package api

import (
	apiserver "github.com/dm0275/mcrun/pkg/api"
	"github.com/spf13/cobra"
)

func NewAPICmd() *cobra.Command {
	var (
		host string
		port int
	)

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the mcrun HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := apiserver.NewServer(host, port)
			cmd.Printf("Starting API server on %s\n", server.Addr())
			return server.Start(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host interface to bind the API server")
	cmd.Flags().IntVar(&port, "port", 8080, "Port that the API server listens on")

	return cmd
}
