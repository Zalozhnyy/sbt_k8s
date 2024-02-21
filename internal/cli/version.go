package cli

import (
	"os"

	"github.com/Zalozhnyy/sbt_k8s/internal/version"
	"github.com/spf13/cobra"
)

func NewVersionCmd() (*cobra.Command, error) {

	return &cobra.Command{
		Use:   "version",
		Short: "Prints the build version",
		Run: func(*cobra.Command, []string) {
			version.Write(os.Stdout)
		},
	}, nil

}
