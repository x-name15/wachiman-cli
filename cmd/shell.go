package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

var shellType string

var ShellCmd = &cobra.Command{
	Use:   "shell [nombre]",
	Short: "Abre una terminal interactiva en el contenedor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		container := args[0]
		
		fmt.Printf("Conectando a %s usando %s...\n", container, shellType)

		c := exec.Command("docker", "exec", "-it", container, shellType)
		
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return fmt.Errorf("error al ejecutar shell (¿el contenedor está corriendo?): %w", err)
		}
		
		return nil
	},
}

func init() {
	ShellCmd.Flags().StringVarP(&shellType, "shell", "s", "/bin/sh", "Shell a usar (/bin/bash, /bin/sh, etc)")
}