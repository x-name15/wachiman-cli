package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Gestiona las redes de Docker",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var networkLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lista todas las redes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		networks, err := client.ListNetworks()
		if err != nil {
			return err
		}

		if len(networks) == 0 {
			fmt.Println("No hay redes.")
			return nil
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, bold("NETWORK ID\tNAME\tDRIVER\tSCOPE"))

		for _, n := range networks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				cyan(n.ID),
				green(n.Name),
				n.Driver,
				n.Scope,
			)
		}
		w.Flush()
		return nil
	},
}

var networkInspectCmd = &cobra.Command{
	Use:   "inspect [nombre]",
	Short: "Muestra detalles de una red",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		n, err := client.InspectNetwork(args[0])
		if err != nil {
			return err
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		internal := "no"
		if n.Internal {
			internal = yellow("sí")
		}

		fmt.Printf("\n%s %s\n\n", bold("⣿ Red:"), cyan(n.Name))
		fmt.Printf("  %s %s\n", bold("ID:      "), n.ID)
		fmt.Printf("  %s %s\n", bold("Driver:  "), n.Driver)
		fmt.Printf("  %s %s\n", bold("Scope:   "), n.Scope)
		fmt.Printf("  %s %s\n", bold("Subnet:  "), n.Subnet)
		fmt.Printf("  %s %s\n", bold("Gateway: "), n.Gateway)
		fmt.Printf("  %s %s\n\n", bold("Internal:"), internal)

		fmt.Printf("%s\n", bold("Contenedores conectados:"))
		if len(n.Containers) == 0 {
			fmt.Println("  ninguno")
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "  "+bold("NAME\tIP"))
			for _, c := range n.Containers {
				fmt.Fprintf(w, "  %s\t%s\n", green(c.Name), cyan(c.IP))
			}
			w.Flush()
		}
		fmt.Println()
		return nil
	},
}

var networkDisconnectCmd = &cobra.Command{
	Use:   "disconnect [red] [contenedor]",
	Short: "Desconecta un contenedor de una red",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		network := args[0]
		container := args[1]

		client, err := docker.New()
		if err != nil {
			return err
		}

		if err := client.NetworkDisconnect(network, container); err != nil {
			return err
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s %s desconectado de %s\n", green("✓"), container, network)
		return nil
	},
}

var networkConnectCmd = &cobra.Command{
	Use:   "connect [red] [contenedor]",
	Short: "Conecta un contenedor a una red",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		network := args[0]
		container := args[1]

		client, err := docker.New()
		if err != nil {
			return err
		}

		if err := client.NetworkConnect(network, container); err != nil {
			return err
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s %s conectado a %s\n", green("✓"), container, network)
		return nil
	},
}

func init() {
	NetworkCmd.AddCommand(networkLsCmd)
	NetworkCmd.AddCommand(networkInspectCmd)
	NetworkCmd.AddCommand(networkConnectCmd)
	NetworkCmd.AddCommand(networkDisconnectCmd)
}