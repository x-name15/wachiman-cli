package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var filterRunning bool
var filterStopped bool
var outputFormat string

var PsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lista todos los contenedores",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		containers, err := client.List()
		if err != nil {
			return err
		}

		filtered := make([]docker.Container, 0)
		for _, c := range containers {
			if filterRunning && !c.Running {
				continue
			}
			if filterStopped && c.Running {
				continue
			}
			filtered = append(filtered, c)
		}
		containers = filtered

		if len(containers) == 0 {
			fmt.Println("No hay contenedores.")
			return nil
		}

		// JSON output
		if outputFormat == "json" {
			data, err := json.MarshalIndent(containers, "", "  ")
			if err != nil {
				return fmt.Errorf("error serializando JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "CONTAINER ID\tNAME\tIMAGE\tSTATUS\tPORTS")

		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		for _, c := range containers {
			status := c.Status
			if c.Running {
				status = green(status)
			} else {
				status = red(status)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", c.ID, c.Name, c.Image, status, c.Ports)
		}
		w.Flush()
		return nil
	},
}

func init() {
	PsCmd.Flags().BoolVar(&filterRunning, "running", false, "Mostrar solo contenedores activos")
	PsCmd.Flags().BoolVar(&filterStopped, "stopped", false, "Mostrar solo contenedores parados")
	PsCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Formato de salida: table, json")
}