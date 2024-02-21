package main

import (
	"log"
	"os"

	"github.com/mdevilliers/go/cli"
	commands "github.com/Zalozhnyy/sbt_k8s/internal/cli"
	"github.com/spf13/cobra"
)

// Application entry point
func main() {
	cmd, err := rootCmd()

	if err != nil {
		log.Fatalf("error configuring commands : %s", err.Error())
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "sbt_k8s",
		Short: "TODO",
	}
	return cmd, cli.RegisterCommands(cmd, commands.NewVersionCmd, commands.NewServerCmd)
}
