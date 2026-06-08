package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var RestartCmd = &cobra.Command{
	Use:   "restart [nombre]",
	Short: "Reinicia un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}
		if err := client.Restart(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ %s reiniciado\n", args[0])
		return nil
	},
}