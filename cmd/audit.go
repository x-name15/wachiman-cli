package cmd

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var auditFix bool

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
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("\n%s %s\n\n", bold("Auditoría del Wachiman:"), cyan(info.Name))

		warnings := 0
		dangers := 0

		if info.User == "root" || info.User == "" {
			fmt.Printf("[ %s ] Usuario: El contenedor corre como %s.\n", red("PELIGRO"), red("root"))
			dangers++
		} else {
			fmt.Printf("[ %s ] Usuario: %s\n", green("OK"), info.User)
		}

		if info.MemoryLimit == 0 {
			fmt.Printf("[ %s ] Memoria: No hay límite configurado — el contenedor puede consumir toda la RAM del host.\n", yellow("WARN"))
			warnings++

			if auditFix {
				fmt.Printf("         %s Aplicando límite por defecto de 512MB...\n", cyan("→"))
				if err := client.SetMemoryLimit(args[0], 512*1024*1024); err != nil {
					fmt.Printf("         %s No se pudo aplicar: %v\n", red("✗"), err)
				} else {
					fmt.Printf("         %s Límite de 512MB aplicado. Reinicia el contenedor para que tome efecto.\n", green("✓"))
				}
			}
		} else {
			mb := float64(info.MemoryLimit) / 1024 / 1024
			fmt.Printf("[ %s ] Memoria: Límite configurado (%.0f MB).\n", green("OK"), mb)
		}

		sensitivePorts := map[string]string{
			"3306":  "MySQL",
			"5432":  "PostgreSQL",
			"6379":  "Redis",
			"27017": "MongoDB",
			"9200":  "Elasticsearch",
			"5984":  "CouchDB",
			"6380":  "Redis (alt)",
		}

		exposedSensitives := false
		warnedPorts := make(map[string]bool)

		for _, p := range info.Ports {
			parts := strings.Split(p, "->")
			if len(parts) == 2 {
				cPort := strings.Split(strings.TrimSpace(parts[1]), "/")[0]
				if dbName, isSensitive := sensitivePorts[cPort]; isSensitive {
					if !warnedPorts[cPort] {
						fmt.Printf("[ %s ] Puerto: %s (%s) expuesto al exterior.\n", yellow("WARN"), yellow(cPort), dbName)
						exposedSensitives = true
						warnedPorts[cPort] = true
						warnings++
					}
				}
			}
		}

		if !exposedSensitives {
			fmt.Printf("[ %s ] Puertos: No se detectaron bases de datos o cachés expuestos.\n", green("OK"))
		}

		dangerousVolume := false
		for _, v := range info.Volumes {
			if strings.Contains(v, "docker.sock") {
				fmt.Printf("[ %s ] Volumen: Acceso al socket de Docker (%s) detectado — riesgo de escape de contenedor.\n", red("PELIGRO"), red("docker.sock"))
				dangerousVolume = true
				dangers++
			} else if strings.HasPrefix(v, "/ ->") {
				fmt.Printf("[ %s ] Volumen: El directorio raíz %s del host está montado en el contenedor.\n", red("PELIGRO"), red("/"))
				dangerousVolume = true
				dangers++
			}
		}

		if !dangerousVolume {
			fmt.Printf("[ %s ] Volumen: No se detectaron montajes críticos.\n", green("OK"))
		}

		if strings.HasSuffix(info.Image, ":latest") || !strings.Contains(info.Image, ":") {
			fmt.Printf("[ %s ] Imagen: Usando tag %s — no garantiza reproducibilidad ni seguridad.\n", yellow("WARN"), yellow("latest"))
			warnings++
		} else {
			fmt.Printf("[ %s ] Imagen: Tag de versión específico detectado (%s).\n", green("OK"), info.Image)
		}

		secretKeywords := []string{"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY", "AUTH", "CREDENTIAL"}
		exposedSecrets := []string{}

		for _, env := range info.Env {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.ToUpper(parts[0])
			val := parts[1]

			for _, keyword := range secretKeywords {
				if strings.Contains(key, keyword) && val != "" && !strings.HasPrefix(val, "${") {
					exposedSecrets = append(exposedSecrets, parts[0])
					break
				}
			}
		}

		if len(exposedSecrets) > 0 {
			fmt.Printf("[ %s ] Secrets: %d variable(s) con datos sensibles en texto plano:\n", red("PELIGRO"), len(exposedSecrets))
			for _, s := range exposedSecrets {
				fmt.Printf("         %s %s\n", red("→"), s)
			}
			dangers++
		} else {
			fmt.Printf("[ %s ] Secrets: No se detectaron secrets expuestos en variables de entorno.\n", green("OK"))
		}

		fmt.Println("\n" + strings.Repeat("─", 50))
		if dangers == 0 && warnings == 0 {
			fmt.Printf("Resultado: %s Contenedor bien configurado.\n", green("¡Excelente!"))
		} else {
			textDangers := fmt.Sprintf("%d peligro", dangers)
			if dangers != 1 {
				textDangers += "s"
			}
			textWarnings := fmt.Sprintf("%d advertencia", warnings)
			if warnings != 1 {
				textWarnings += "s"
			}
			fmt.Printf("Resultado: Se encontraron %s y %s.\n", red(textDangers), yellow(textWarnings))

			if !auditFix && (dangers > 0 || warnings > 0) {
				fmt.Printf("\n%s Usa %s para intentar corregir automáticamente lo que sea posible.\n",
					cyan("→"), cyan("wachiman audit "+args[0]+" --fix"))
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	AuditCmd.Flags().BoolVar(&auditFix, "fix", false, "Intentar corregir automáticamente los problemas detectados")
}