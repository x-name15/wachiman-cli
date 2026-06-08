package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var StopCmd = &cobra.Command{
	Use:   "stop [nombre]",
	Short: "Para un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}
		if err := client.Stop(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ %s parado\n", args[0])
		return nil
	},
}