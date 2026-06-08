package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var watchInterval int

var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor en tiempo real de contenedores",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
		defer ticker.Stop()
		draw(client)

		for {
			select {
			case <-ticker.C:
				draw(client)
			case <-sig:
				fmt.Print("\033[H\033[2J")
				fmt.Println("Wachiman saliendo...")
				return nil
			}
		}
	},
}

func draw(client *docker.Client) {
	fmt.Print("\033[H\033[2J")

	overview, err := client.Overview()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Contar corriendo vs parados
	running := 0
	stopped := 0
	for _, c := range overview {
		if c.Running {
			running++
		} else {
			stopped++
		}
	}

	fmt.Printf("%s %s\n%s\n\n",
		cyan(bold("⣿ wachiman watch")),
		color.New(color.FgHiBlack).Sprintf("— actualizado: %s — Ctrl+C para salir", time.Now().Format("15:04:05")),
		fmt.Sprintf("%s   %s",
			green(fmt.Sprintf("● %d corriendo", running)),
			red(fmt.Sprintf("● %d parado", stopped)),
		),
	)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, bold("CONTAINER ID\tNAME\tIMAGE\tSTATUS\tCPU %\tMEM %\tPORTS"))

	for _, c := range overview {
		status := green(c.Status)
		if !c.Running {
			status = red(c.Status)
		}

		cpuStr := fmt.Sprintf("%s %5.1f%%", renderBar(c.CPU, 10), c.CPU)
		memStr := fmt.Sprintf("%s %5.1f%%", renderBar(c.Memory, 10), c.Memory)

		if !c.Running {
			cpuStr = "-"
			memStr = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			c.ID, c.Name, c.Image, status, cpuStr, memStr, c.Ports)
	}
	w.Flush()
}

func init() {
	WatchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 3, "Intervalo de refresco en segundos")
}

func renderBar(percent float64, width int) string {
	filled := int(percent / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	// Color según porcentaje
	if percent > 80 {
		return color.RedString(bar)
	} else if percent > 50 {
		return color.YellowString(bar)
	}
	return color.GreenString(bar)
}