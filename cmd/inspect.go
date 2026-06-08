package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var InspectCmd = &cobra.Command{
	Use:   "inspect [nombre]",
	Short: "Muestra detalles de un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		info, err := client.Inspect(args[0])
		if err != nil {
			return err
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		status := green(info.Status)
		if info.Status != "running" {
			status = red(info.Status)
		}

		fmt.Printf("\n%s %s\n", bold("Contenedor:"), cyan(info.Name))
		fmt.Printf("%s %s\n", bold("Estado:       "), status)
		fmt.Printf("%s %s\n", bold("Imagen:       "), info.Image)
		fmt.Printf("%s %s\n\n", bold("IP:           "), info.IP)

		fmt.Printf("%s\n", bold("Puertos:"))
		if len(info.Ports) == 0 {
			fmt.Println("  ninguno")
		}
		for _, p := range info.Ports {
			fmt.Printf("  %s\n", p)
		}

		fmt.Printf("\n%s\n", bold("Volúmenes:"))
		if len(info.Volumes) == 0 {
			fmt.Println("  ninguno")
		}
		for _, v := range info.Volumes {
			fmt.Printf("  %s\n", v)
		}

		fmt.Printf("\n%s\n", bold("Variables de entorno:"))
		for _, e := range info.Env {
			fmt.Printf("  %s\n", e)
		}

		fmt.Println()
		return nil
	},
}