package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const Version = "1.9.9"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión de wachiman",
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("%s %s\n", bold("wachiman"), cyan("v"+Version))
	},
}