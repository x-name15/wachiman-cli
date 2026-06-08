package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/config"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Gestiona la configuración de wachiman",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Muestra la configuración actual",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("\n%s\n\n", bold("⣿ wachiman config"))
		fmt.Printf("  %s %v\n", cyan("watch_interval"), fmt.Sprintf("%ds", cfg.WatchInterval))
		fmt.Printf("  %s %v\n", cyan("default_tail  "), cfg.DefaultTail)
		fmt.Printf("  %s %v\n\n", cyan("output_format "), cfg.OutputFormat)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Modifica un valor de configuración",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Set(args[0], args[1]); err != nil {
			return err
		}
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s %s = %s\n", green("✓"), args[0], args[1])
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Restaura la configuración a los valores por defecto",
	RunE: func(cmd *cobra.Command, args []string) error {
		defaults := &config.Config{
			WatchInterval: 3,
			DefaultTail:   50,
			OutputFormat:  "table",
		}
		if err := config.Save(defaults); err != nil {
			return err
		}
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("✓ configuración restaurada a defaults"))
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configResetCmd)
}