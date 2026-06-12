package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var DiffCmd = &cobra.Command{
	Use:   "diff [nombre]",
	Short: "Muestra los archivos modificados dentro de un contenedor desde que arrancó",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		changes, err := client.Diff(args[0])
		if err != nil {
			return err
		}

		if len(changes) == 0 {
			fmt.Println("Sin cambios desde el arranque.")
			return nil
		}

		bold := color.New(color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("\n%s %s\n\n", bold("⣿ diff:"), cyan(args[0]))
		fmt.Printf("  %s añadido   %s modificado   %s eliminado\n\n",
			green("A"),
			yellow("M"),
			red("D"),
		)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

		added, modified, deleted := 0, 0, 0

		for _, ch := range changes {
			switch ch.Kind {
			case "A":
				fmt.Fprintf(w, "  %s\t%s\n", green("A"), ch.Path)
				added++
			case "M":
				fmt.Fprintf(w, "  %s\t%s\n", yellow("M"), ch.Path)
				modified++
			case "D":
				fmt.Fprintf(w, "  %s\t%s\n", red("D"), ch.Path)
				deleted++
			}
		}
		w.Flush()

		fmt.Printf("\n%s %s añadidos   %s modificados   %s eliminados\n\n",
			bold("→"),
			green(fmt.Sprintf("%d", added)),
			yellow(fmt.Sprintf("%d", modified)),
			red(fmt.Sprintf("%d", deleted)),
		)

		return nil
	},
}