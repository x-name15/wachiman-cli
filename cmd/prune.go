package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var forceprune bool

var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Limpia contenedores parados, imágenes y volúmenes sin usar",
	RunE: func(cmd *cobra.Command, args []string) error {
		yellow := color.New(color.FgYellow).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		
		if !forceprune {
			fmt.Printf("%s Esto eliminará contenedores parados, imágenes sin usar y volúmenes huérfanos.\n", yellow("⚠"))
			fmt.Printf("%s ¿Continuar? [y/N]: ", bold("→"))

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input != "y" && input != "yes" {
				fmt.Println("Cancelado.")
				return nil
			}
		}

		client, err := docker.New()
		if err != nil {
			return err
		}

		fmt.Println("Limpiando...")
		summary, err := client.Prune()
		if err != nil {
			return err
		}

		mb := float64(summary.SpaceReclaimed) / 1024 / 1024
		fmt.Printf("\n%s\n", bold("Resultado:"))
		fmt.Printf("  %s contenedores eliminados\n", green(fmt.Sprintf("%d", summary.ContainersDeleted)))
		fmt.Printf("  %s imágenes eliminadas\n", green(fmt.Sprintf("%d", summary.ImagesDeleted)))
		fmt.Printf("  %s volúmenes eliminados\n", green(fmt.Sprintf("%d", summary.VolumesDeleted)))
		fmt.Printf("  %s liberados\n", green(fmt.Sprintf("%.2f MB", mb)))
		return nil
	},
}

func init() {
	PruneCmd.Flags().BoolVarP(&forceprune, "force", "f", false, "Sin confirmación")
}