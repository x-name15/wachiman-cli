package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
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

		fmt.Printf("%s Inspeccionando contenedor %s...\n", bold("→"), cyan(containerName))
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
					fmt.Printf("%s Error al reanudar: %v\n", red("❌"), err)
				}
			}()
		} else if info.Status == "running" && noPause {
			fmt.Printf("%s %s Ejecutando respaldo en caliente. La data activa podría ser inconsistente.\n", yellow("⚠"), bold("CUIDADO:"))
		}

		if err := os.MkdirAll(backupOutputDir, 0755); err != nil {
			return fmt.Errorf("no se pudo crear la carpeta de salida: %w", err)
		}

		timestamp := time.Now().Format("20060102_150405")

		if len(info.Volumes) > 0 {
			fmt.Printf("Se detectaron %d volúmenes montados. Analizando rutas...\n", len(info.Volumes))

			var rawPaths []string
			for _, v := range info.Volumes {
				parts := strings.Split(v, " -> ")
				if len(parts) == 2 {
					rawPaths = append(rawPaths, parts[1])
				}
			}

			var targetPaths []string
			for _, p := range rawPaths {
				isSubPath := false
				for _, other := range rawPaths {
					if p != other && strings.HasPrefix(p, other+"/") {
						isSubPath = true
						break
					}
				}
				if !isSubPath {
					targetPaths = append(targetPaths, p)
				}
			}

			for i, containerPath := range targetPaths {
				safePathName := strings.ReplaceAll(strings.Trim(containerPath, "/"), "/", "-")
				if safePathName == "" {
					safePathName = fmt.Sprintf("vol-%d", i)
				}

				fileName := fmt.Sprintf("%s_%s_%s.tar", containerName, safePathName, timestamp)
				outputPath := filepath.Join(backupOutputDir, fileName)

				fmt.Printf("\n%s Respaldando %s\n", bold("→"), cyan(containerPath))

				reader, err := client.CopyFrom(containerName, containerPath)
				if err != nil {
					return fmt.Errorf("%s falló: %w", containerPath, err)
				}

				size, err := saveBackupFileWithProgress(reader, outputPath, fileName)
				reader.Close()

				if err != nil {
					return fmt.Errorf("error guardando %s: %w", fileName, err)
				}

				fmt.Printf("%s Guardado: %s (%s)\n", green("✓"), cyan(fileName), formatSize(size))
			}
		} else {
			fmt.Printf("%s No se encontraron volúmenes. Exportando el sistema raíz (rootfs)...\n", cyan("→"))

			fileName := fmt.Sprintf("%s_rootfs_%s.tar", containerName, timestamp)
			outputPath := filepath.Join(backupOutputDir, fileName)

			fmt.Printf("\n%s Exportando %s\n", bold("→"), yellow(fileName))

			reader, err := client.Export(containerName)
			if err != nil {
				return fmt.Errorf("export falló: %w", err)
			}

			size, err := saveBackupFileWithProgress(reader, outputPath, fileName)
			reader.Close()

			if err != nil {
				return fmt.Errorf("error guardando %s: %w", fileName, err)
			}

			fmt.Printf("%s Guardado: %s (%s)\n", green("✓"), cyan(fileName), formatSize(size))
		}

		fmt.Printf("\n%s %s Respaldo completado en: %s\n", green("✓"), bold("Wachiman dice:"), cyan(backupOutputDir))
		return nil
	},
}

func saveBackupFileWithProgress(reader io.Reader, path, label string) (int64, error) {
	outFile, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()
	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription(label),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	written, err := io.Copy(io.MultiWriter(outFile, bar), reader)
	bar.Finish()
	return written, err
}

func formatSize(b int64) string {
	mb := float64(b) / 1024 / 1024
	if mb >= 1024 {
		return fmt.Sprintf("%.2f GB", mb/1024)
	}
	return fmt.Sprintf("%.2f MB", mb)
}

func init() {
	BackupCmd.Flags().StringVarP(&backupOutputDir, "output", "o", ".", "Directorio donde se guardará el archivo .tar")
	BackupCmd.Flags().BoolVar(&noPause, "no-pause", false, "Hacer respaldo en vivo sin pausar el contenedor")
}