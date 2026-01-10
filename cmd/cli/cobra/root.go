package cobra

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "service",
	Short: "service CLI",
}

func GetRootCmd() *cobra.Command { return rootCmd }

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
