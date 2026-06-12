package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/config"
	"wachiman/docker"
)

var watchInterval int
var sparks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

type history struct {
	CPU []float64
	Mem []float64
}

var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor en tiempo real de contenedores",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if !cmd.Flags().Changed("interval") {
			watchInterval = cfg.WatchInterval
		}

		client, err := docker.New()
		if err != nil {
			return err
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
		defer ticker.Stop()

		hist := make(map[string]*history)
		draw(client, hist)

		for {
			select {
			case <-ticker.C:
				draw(client, hist)
			case <-sig:
				clearScreen()
				fmt.Println("wachiman saliendo...")
				return nil
			}
		}
	},
}

func draw(client *docker.Client, hist map[string]*history) {
	fmt.Print(clearScreen())

	overview, err := client.Overview()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

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
	fmt.Fprintln(w, bold("CONTAINER ID\tNAME\tIMAGE\tSTATUS\tCPU\tMEM\tPORTS"))

	for _, c := range overview {
		status := green(c.Status)
		if !c.Running {
			status = red(c.Status)
		}

		if _, ok := hist[c.Name]; !ok {
			hist[c.Name] = &history{}
		}
		h := hist[c.Name]

		if c.Running {
			h.CPU = appendCapped(h.CPU, c.CPU, 10)
			h.Mem = appendCapped(h.Mem, c.Memory, 10)
		}

		cpuStr := "-"
		memStr := "-"
		if c.Running {
			cpuBar := renderBar(c.CPU, 8)
			memBar := renderBar(c.Memory, 8)
			cpuSpark := renderSparkline(h.CPU)
			memSpark := renderSparkline(h.Mem)
			cpuStr = fmt.Sprintf("%s %s %4.1f%%", cpuSpark, cpuBar, c.CPU)
			memStr = fmt.Sprintf("%s %s %4.1f%%", memSpark, memBar, c.Memory)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			c.ID, c.Name, c.Image, status, cpuStr, memStr, c.Ports)
	}
	w.Flush()
}

func appendCapped(slice []float64, val float64, maxLen int) []float64 {
	slice = append(slice, val)
	if len(slice) > maxLen {
		slice = slice[len(slice)-maxLen:]
	}
	return slice
}

func renderSparkline(values []float64) string {
	if len(values) == 0 {
		return ""
	}

	result := ""
	for _, v := range values {
		idx := int(v / 100.0 * float64(len(sparks)-1))
		if idx >= len(sparks) {
			idx = len(sparks) - 1
		}
		if idx < 0 {
			idx = 0
		}
		result += string(sparks[idx])
	}

	last := values[len(values)-1]
	if last > 80 {
		return color.RedString(result)
	} else if last > 50 {
		return color.YellowString(result)
	}
	return color.GreenString(result)
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

	if percent > 80 {
		return color.RedString(bar)
	} else if percent > 50 {
		return color.YellowString(bar)
	}
	return color.GreenString(bar)
}

func init() {
	WatchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 3, "Intervalo de refresco en segundos")
}

func clearScreen() string {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return ""
}