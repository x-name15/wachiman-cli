package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var TopCmd = &cobra.Command{
	Use:   "top [nombre]",
	Short: "Muestra los procesos corriendo dentro de un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		processes, err := client.Top(args[0])
		if err != nil {
			return err
		}

		if len(processes) == 0 {
			fmt.Println("No hay procesos.")
			return nil
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("\n%s %s\n\n", bold("⣿ procesos en"), cyan(args[0]))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, bold("PID\tCOMMAND\tCPU %\tMEM %"))

		for _, p := range processes {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.PID, p.Name, p.CPU, p.Mem)
		}
		w.Flush()
		fmt.Println()
		return nil
	},
}