package cmd

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var ExportComposeCmd = &cobra.Command{
	Use:   "export-compose [nombre]",
	Short: "Exporta la configuración de un contenedor a formato docker-compose.yml",
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

		cyan := color.New(color.FgCyan).SprintFunc()
		
		var sb strings.Builder

		sb.WriteString("services:\n")
		sb.WriteString(fmt.Sprintf("  %s:\n", info.Name))
		sb.WriteString(fmt.Sprintf("    image: %s\n", info.Image))

		if len(info.Ports) > 0 {
			sb.WriteString("    ports:\n")
			seenPorts := make(map[string]bool)
			for _, p := range info.Ports {
				parts := strings.Split(p, "->")
				if len(parts) == 2 {
					hostPart := parts[0]
					if idx := strings.LastIndex(hostPart, ":"); idx != -1 {
						hostPart = hostPart[idx+1:]
					}

					containerPart := strings.Split(parts[1], "/")[0]
					
					mapping := fmt.Sprintf("%s:%s", hostPart, containerPart)
					if !seenPorts[mapping] {
						sb.WriteString(fmt.Sprintf("      - \"%s\"\n", mapping))
						seenPorts[mapping] = true
					}
				}
			}
		}

		if len(info.Env) > 0 {
			sb.WriteString("    environment:\n")
			for _, env := range info.Env {
				if !strings.HasPrefix(env, "PATH=") && !strings.HasPrefix(env, "HOME=") {
					sb.WriteString(fmt.Sprintf("      - %s\n", env))
				}
			}
		}

		if len(info.Volumes) > 0 {
			sb.WriteString("    volumes:\n")
			for _, v := range info.Volumes {
				parts := strings.Split(v, " -> ")
				if len(parts) == 2 {
					sb.WriteString(fmt.Sprintf("      - %s:%s\n", parts[0], parts[1]))
				} else {
					sb.WriteString(fmt.Sprintf("      - %s\n", v))
				}
			}
		}

		sb.WriteString("    restart: unless-stopped\n")

		fmt.Printf("\n%s\n\n", cyan("# Generado por Wachiman CLI"))
		fmt.Println(sb.String())

		return nil
	},
}