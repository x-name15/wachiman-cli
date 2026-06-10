package cmd

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var AuditCmd = &cobra.Command{
	Use:   "audit [nombre]",
	Short: "Audita la configuración de seguridad de un contenedor",
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
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("\n%s %s\n\n", bold("Auditoría del Wachiman:"), info.Name)

		warnings := 0
		dangers := 0

		if info.User == "root" || info.User == "root (implícito)" || info.User == "" {
			fmt.Printf("[ %s ] Usuario: El contenedor corre como root.\n", red("PELIGRO"))
			dangers++
		} else {
			fmt.Printf("[ %s ] Usuario: %s\n", green("OK"), info.User)
		}

		if info.MemoryLimit == 0 {
			fmt.Printf("[ %s ] Memoria: No hay límite de memoria configurado.\n", yellow("WARN"))
			warnings++
		} else {
			mb := float64(info.MemoryLimit) / 1024 / 1024
			fmt.Printf("[ %s ] Memoria: Límite configurado (%.2f MB).\n", green("OK"), mb)
		}

		sensitivePorts := map[string]string{
			"3306":  "MySQL",
			"5432":  "PostgreSQL",
			"6379":  "Redis",
			"27017": "MongoDB",
			"9200":  "Elasticsearch",
		}

		exposedSensitives := false
		warnedPorts := make(map[string]bool)

		for _, p := range info.Ports {
			parts := strings.Split(p, "->")
			if len(parts) == 2 {
				cPort := strings.Split(parts[1], "/")[0]
				
				if dbName, isSensitive := sensitivePorts[cPort]; isSensitive {
					if !warnedPorts[cPort] {
						fmt.Printf("[ %s ] Puerto: %s (%s) expuesto hacia afuera.\n", yellow("WARN"), cPort, dbName)
						exposedSensitives = true
						warnedPorts[cPort] = true
						warnings++
					}
				}
			}
		}

		if len(info.Ports) == 0 || !exposedSensitives {
			fmt.Printf("[ %s ] Puertos: No se detectaron bases de datos o cachés expuestos.\n", green("OK"))
		}
		dangerousVolume := false
		for _, v := range info.Volumes {
			if strings.Contains(v, "docker.sock") {
				fmt.Printf("[ %s ] Volumen: ¡Acceso al socket de Docker (docker.sock) detectado!\n", red("PELIGRO"))
				dangerousVolume = true
				dangers++
			} else if strings.HasPrefix(v, "/ ->") {
				fmt.Printf("[ %s ] Volumen: El directorio raíz (/) del host está expuesto al contenedor.\n", red("PELIGRO"))
				dangerousVolume = true
				dangers++
			}
		}

		if !dangerousVolume {
			fmt.Printf("[ %s ] Volumen: No se detectaron montajes críticos.\n", green("OK"))
		}

		fmt.Println("\n" + strings.Repeat("-", 40))
		if dangers == 0 && warnings == 0 {
			fmt.Printf("Resultado: %s Contenedor bien configurado.\n", green("¡Excelente!"))
		} else {
			textDangers := fmt.Sprintf("%d peligro", dangers)
			if dangers != 1 { textDangers += "s" }

			textWarnings := fmt.Sprintf("%d advertencia", warnings)
			if warnings != 1 { textWarnings += "s" }

			fmt.Printf("Resultado: Se encontraron %s y %s.\n", red(textDangers), yellow(textWarnings))
		}
		fmt.Println()

		return nil
	},
}