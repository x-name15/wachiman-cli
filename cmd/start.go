package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var StartCmd = &cobra.Command{
	Use:   "start [nombre]",
	Short: "Arranca un contenedor parado",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}
		if err := client.Start(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ %s arrancado\n", args[0])
		return nil
	},
}