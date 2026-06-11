package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var (
	monitorInterval  int
	monitorNoRestart bool
	monitorOnly      string
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitorea contenedores y los reinicia si se caen",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		filter := []string{}
		if monitorOnly != "" {
			for _, name := range strings.Split(monitorOnly, ",") {
				filter = append(filter, strings.TrimSpace(name))
			}
		}

		fmt.Printf("\n%s\n", bold("⣿ wachiman monitor"))
		if monitorNoRestart {
			fmt.Printf("%s Modo observación — no se reiniciará ningún contenedor\n", yellow("👁"))
		} else {
			fmt.Printf("%s Modo activo — contenedores caídos serán reiniciados automáticamente\n", cyan("●"))
		}
		if len(filter) > 0 {
			fmt.Printf("%s Monitoreando: %s\n", bold("→"), strings.Join(filter, ", "))
		} else {
			fmt.Printf("%s Monitoreando todos los contenedores\n", bold("→"))
		}
		fmt.Printf("%s Intervalo: %ds — Ctrl+C para salir\n\n", bold("→"), monitorInterval)

		prevState := make(map[string]string)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(monitorInterval) * time.Second)
		defer ticker.Stop()

		initMonitor(client, prevState, filter)
		fmt.Printf("%s Estado inicial capturado. Vigilando...\n\n", cyan("✓"))

		for {
			select {
			case <-ticker.C:
				runMonitor(client, prevState, filter, monitorNoRestart)
			case <-sig:
				fmt.Printf("\n%s wachiman monitor detenido.\n", bold("→"))
				return nil
			}
		}
	},
}

func simpleState(running bool) string {
	if running {
		return "running"
	}
	return "stopped"
}

func initMonitor(client *docker.Client, prevState map[string]string, filter []string) {
	containers, err := client.List()
	if err != nil {
		return
	}

	for _, c := range containers {
		if !shouldMonitor(c.Name, filter) {
			continue
		}
		prevState[c.Name] = simpleState(c.Running)
	}
}

func runMonitor(client *docker.Client, prevState map[string]string, filter []string, noRestart bool) {
	containers, err := client.List()
	if err != nil {
		fmt.Printf("error listando contenedores: %v\n", err)
		return
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	timestamp := time.Now().Format("15:04:05")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	for _, c := range containers {
		if !shouldMonitor(c.Name, filter) {
			continue
		}

		current := simpleState(c.Running)
		prev, exists := prevState[c.Name]

		if !exists {
			prevState[c.Name] = current
			fmt.Fprintf(w, "%s\t%s %s\t%s\n",
				cyan(timestamp),
				"●",
				bold(c.Name),
				"nuevo contenedor detectado",
			)
			w.Flush()
			continue
		}

		// Comparamos estado simplificado — no el string completo
		if prev == current {
			continue
		}

		prevState[c.Name] = current

		if c.Running {
			fmt.Fprintf(w, "%s\t%s %s\t%s\n",
				cyan(timestamp),
				green("●"),
				bold(c.Name),
				green("volvió a estar activo"),
			)
			w.Flush()
		} else {
			fmt.Fprintf(w, "%s\t%s %s\t%s\n",
				cyan(timestamp),
				red("●"),
				bold(c.Name),
				red(c.Status),
			)
			w.Flush()

			if !noRestart {
				fmt.Fprintf(w, "%s\t  %s Reiniciando %s...\n",
					cyan(timestamp),
					yellow("↺"),
					bold(c.Name),
				)
				w.Flush()

				if err := client.Restart(c.Name); err != nil {
					fmt.Fprintf(w, "%s\t  %s Error al reiniciar: %v\n",
						cyan(timestamp),
						red("✗"),
						err,
					)
				} else {
					fmt.Fprintf(w, "%s\t  %s %s reiniciado correctamente\n",
						cyan(timestamp),
						green("✓"),
						bold(c.Name),
					)
				}
				w.Flush()
			}
		}
	}
}

func shouldMonitor(name string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	for _, f := range filter {
		if f == name {
			return true
		}
	}
	return false
}

func init() {
	MonitorCmd.Flags().IntVarP(&monitorInterval, "interval", "i", 5, "Intervalo de chequeo en segundos")
	MonitorCmd.Flags().BoolVar(&monitorNoRestart, "no-restart", false, "Solo observar, no reiniciar")
	MonitorCmd.Flags().StringVar(&monitorOnly, "only", "", "Contenedores a monitorear (separados por coma)")
}