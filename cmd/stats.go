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

var statsOutputFormat string

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Muestra CPU y memoria de los contenedores activos",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		stats, err := client.Stats()
		if err != nil {
			return err
		}

		if len(stats) == 0 {
			fmt.Println("No hay contenedores corriendo.")
			return nil
		}

		if statsOutputFormat == "json" {
			data, err := json.MarshalIndent(stats, "", "  ")
			if err != nil {
				return fmt.Errorf("error serializando JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tCPU %\tMEM USED\tMEM LIMIT\tMEM %")

		for _, s := range stats {
			cpuStr := fmt.Sprintf("%.2f%%", s.CPU)
			memStr := fmt.Sprintf("%.2f%%", s.Memory)

			if s.CPU > 80 {
				cpuStr = color.RedString(cpuStr)
			} else if s.CPU > 50 {
				cpuStr = color.YellowString(cpuStr)
			} else {
				cpuStr = color.GreenString(cpuStr)
			}

			if s.Memory > 80 {
				memStr = color.RedString(memStr)
			} else if s.Memory > 50 {
				memStr = color.YellowString(memStr)
			} else {
				memStr = color.GreenString(memStr)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.Name, cpuStr, s.MemUsed, s.MemLimit, memStr)
		}
		w.Flush()
		return nil
	},
}

func init() {
	StatsCmd.Flags().StringVarP(&statsOutputFormat, "output", "o", "table", "Formato de salida: table, json")
}