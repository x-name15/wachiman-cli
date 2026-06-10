package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/cmd"
)

const banner = "" +
" __          __        _     _                        \n" +
" \\ \\        / /       | |   (_)                       \n" +
"  \\ \\  /\\  / /_ _  ___| |__  _ _ __ ___   __ _ _ __   \n" +
"   \\ \\/  \\/ / _` |/ __| '_ \\| | '_ ` _ \\ / _` | '_ \\  \n" +
"    \\  /\\  / (_| | (__| | | | | | | | | | (_| | | | | \n" +
"     \\/  \\/ \\__,_|\\___|_| |_|_|_| |_| |_|\\__,_|_| |_| \n"


func main() {
	root := &cobra.Command{
		Use:   "wachiman",
		Short: "Desde Perú hasta tu deploy de Docker, wachiman te cuida",
		Long:  "wachiman — vigilancia y control de contenedores Docker desde tu terminal.",
		// Esto se ejecuta cuando corres "wachiman" sin argumentos
		Run: func(cmd *cobra.Command, args []string) {
			cyan := color.New(color.FgCyan).SprintFunc()
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Println(cyan(banner))
			fmt.Println(yellow("  Desde Perú hasta tu deploy de Docker, wachiman te cuida\n"))
			cmd.Help()
		},
	}

	root.AddCommand(cmd.PsCmd)
	root.AddCommand(cmd.StopCmd)
	root.AddCommand(cmd.RestartCmd)
	root.AddCommand(cmd.LogsCmd)
	root.AddCommand(cmd.StartCmd)
	root.AddCommand(cmd.StatsCmd)
	root.AddCommand(cmd.InspectCmd)
	root.AddCommand(cmd.WatchCmd)
	root.AddCommand(cmd.TopCmd)
	root.AddCommand(cmd.PruneCmd)
	root.AddCommand(cmd.ConfigCmd)
	root.AddCommand(cmd.AuditCmd)
	root.AddCommand(cmd.ShellCmd)
	root.AddCommand(cmd.ExportComposeCmd)
	root.AddCommand(cmd.BackupCmd)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}