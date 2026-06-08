package cmd

import (
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var tail string

var LogsCmd = &cobra.Command{
	Use:   "logs [nombre]",
	Short: "Muestra los logs de un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}
		return client.Logs(args[0], tail)
	},
}

func init() {
	LogsCmd.Flags().StringVarP(&tail, "tail", "t", "50", "Número de líneas a mostrar")
}