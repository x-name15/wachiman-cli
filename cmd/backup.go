package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var (
	backupOutputDir string
	noPause         bool
)

var BackupCmd = &cobra.Command{
	Use:   "backup [nombre]",
	Short: "Crea un respaldo seguro (.tar) de la data de un contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		containerName := args[0]

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		client, err := docker.New()
		if err != nil {
			return err
		}
		fmt.Printf("%s Inspeccionando contenedor %s...\n", bold(containerName), cyan(containerName))
		info, err := client.Inspect(containerName)
		if err != nil {
			return err
		}

		shouldPause := info.Status == "running" && !noPause

		if shouldPause {
			fmt.Printf("%s Pausando contenedor para evitar corrupción de archivos...\n", yellow("⏸"))
			if err := client.Pause(containerName); err != nil {
				return err
			}

			defer func() {
				fmt.Printf("%s Reanudando contenedor...\n", green("▶"))
				if err := client.Unpause(containerName); err != nil {
					fmt.Printf("%s Error al reanudar de inmediato: %v\n", red("❌"), err)
				}
			}()
		} else if info.Status == "running" && noPause {
			fmt.Printf("%s %s Ejecutando respaldo en caliente (hot backup). La data activa podría ser inconsistente.\n", yellow("⚠"), bold("CUIDADO:"))
		}

		if err := os.MkdirAll(backupOutputDir, 0755); err != nil {
			return fmt.Errorf("no se pudo crear la carpeta de salida: %w", err)
		}

		timestamp := time.Now().Format("20060102_150405")

		if len(info.Volumes) > 0 {
			fmt.Printf("Se detectaron %d volúmenes montados. Analizando rutas...\n", len(info.Volumes))

			// 1. Extraer y limpiar todas las rutas válidas del contenedor
			var rawPaths []string
			for _, v := range info.Volumes {
				parts := strings.Split(v, " -> ")
				if len(parts) == 2 {
					rawPaths = append(rawPaths, parts[1])
				}
			}

			// 2. Filtrar subrutas redundantes (ej: no respaldar /a/b si ya se respalda /a)
			var targetPaths []string
			for _, p := range rawPaths {
				isSubPath := false
				for _, other := range rawPaths {
					// Aseguramos que termine con "/" para evitar falsos positivos como /var/www-data contra /var/www
					if p != other && strings.HasPrefix(p, other+"/") {
						isSubPath = true
						break
					}
				}
				if !isSubPath {
					targetPaths = append(targetPaths, p)
				}
			}

			// 3. Procesar solo las rutas optimizadas
			for i, containerPath := range targetPaths {
				safePathName := strings.ReplaceAll(strings.Trim(containerPath, "/"), "/", "-")
				if safePathName == "" {
					safePathName = fmt.Sprintf("vol-%d", i)
				}

				fileName := fmt.Sprintf("%s_%s_%s.tar", containerName, safePathName, timestamp)
				outputPath := filepath.Join(backupOutputDir, fileName)

				fmt.Printf("   %s Respaldando %s %s ", bold("→"), cyan(containerPath), bold("..."))

				reader, err := client.CopyFrom(containerName, containerPath)
				if err != nil {
					fmt.Println(red("¡FALLÓ!"))
					return err
				}

				err = saveBackupFile(reader, outputPath)
				reader.Close()

				if err != nil {
					fmt.Println(red("¡FALLÓ AL GUARDAR!"))
					return err
				}
				fmt.Println(green("✓ Guardado como " + fileName))
			}
		} else {
			fmt.Printf("%s No se encontraron volúmenes. Exportando el sistema raíz (rootfs)...\n", cyan("CD"))

			fileName := fmt.Sprintf("%s_rootfs_%s.tar", containerName, timestamp)
			outputPath := filepath.Join(backupOutputDir, fileName)

			fmt.Printf("   %s Exportando a %s %s ", bold("→"), yellow(fileName), bold("..."))

			reader, err := client.Export(containerName)
			if err != nil {
				fmt.Println(red("¡FALLÓ!"))
				return err
			}

			err = saveBackupFile(reader, outputPath)
			reader.Close()

			if err != nil {
				fmt.Println(red("¡FALLÓ AL GUARDAR!"))
				return err
			}
			fmt.Println(green("✓ Hecho"))
		}

		fmt.Printf("\n%s %s ¡Respaldo completado con éxito en: %s!\n", green("✓"), bold("Wachiman dice:"), cyan(backupOutputDir))
		return nil
	},
}

func saveBackupFile(reader io.Reader, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}

func init() {
	BackupCmd.Flags().StringVarP(&backupOutputDir, "output", "o", ".", "Directorio local donde se guardará el archivo .tar")
	BackupCmd.Flags().BoolVar(&noPause, "no-pause", false, "Hacer respaldo en vivo sin pausar el contenedor")
}