package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ComposeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Gestiona docker compose desde el directorio actual",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func runCompose(args ...string) error {
	found := false
	for _, name := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"} {
		if _, err := os.Stat(name); err == nil {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no se encontró docker-compose.yml en el directorio actual")
	}

	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		cmd2 := exec.Command("docker-compose", args...)
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		cmd2.Stdin = os.Stdin
		return cmd2.Run()
	}
	return nil
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Arranca los servicios del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		detach, _ := cmd.Flags().GetBool("detach")
		composeArgs := []string{"up"}
		if detach {
			composeArgs = append(composeArgs, "-d")
		}
		return runCompose(composeArgs...)
	},
}

var composeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Para y elimina los servicios del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		volumes, _ := cmd.Flags().GetBool("volumes")
		composeArgs := []string{"down"}
		if volumes {
			composeArgs = append(composeArgs, "-v")
		}
		return runCompose(composeArgs...)
	},
}

var composePsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lista los servicios del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCompose("ps")
	},
}

var composeLogsCmd = &cobra.Command{
	Use:   "logs [servicio]",
	Short: "Muestra logs de los servicios del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		tail, _ := cmd.Flags().GetString("tail")
		composeArgs := []string{"logs", "--tail", tail}
		if len(args) > 0 {
			composeArgs = append(composeArgs, args[0])
		}
		return runCompose(composeArgs...)
	},
}

var composeRestartCmd = &cobra.Command{
	Use:   "restart [servicio]",
	Short: "Reinicia los servicios del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeArgs := []string{"restart"}
		if len(args) > 0 {
			composeArgs = append(composeArgs, args[0])
		}
		return runCompose(composeArgs...)
	},
}

var composePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pullea las imágenes del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCompose("pull")
	},
}

var composeBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Buildea las imágenes del compose",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCompose("build")
	},
}

func init() {
	composeUpCmd.Flags().BoolP("detach", "d", false, "Arrancar en background")
	composeDownCmd.Flags().BoolP("volumes", "v", false, "Eliminar volúmenes también")
	composeLogsCmd.Flags().StringP("tail", "t", "50", "Número de líneas a mostrar")

	ComposeCmd.AddCommand(composeUpCmd)
	ComposeCmd.AddCommand(composeDownCmd)
	ComposeCmd.AddCommand(composePsCmd)
	ComposeCmd.AddCommand(composeLogsCmd)
	ComposeCmd.AddCommand(composeRestartCmd)
	ComposeCmd.AddCommand(composePullCmd)
	ComposeCmd.AddCommand(composeBuildCmd)

	ComposeCmd.Run = func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("\n%s\n\n", bold("⣿ wachiman compose"))

		// Detectar archivo compose
		for _, name := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"} {
			if _, err := os.Stat(name); err == nil {
				fmt.Printf("Archivo detectado: %s\n\n", cyan(name))
				break
			}
		}
		cmd.Help()
	}
}