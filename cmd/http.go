package cmd

import (
	"log"
	"nodnotes-api/app"

	"github.com/spf13/cobra"
)

func Http(cmd *cobra.Command, args []string) {
	err := app.Start()
	if err != nil {
		log.Fatalf("api start failed: %v\n", err)
	}
}

func init() {
	cmd := &cobra.Command{
		Use: "http",
		Run: Http,
	}
	rootCmd.AddCommand(cmd)
}
