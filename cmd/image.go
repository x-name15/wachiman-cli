package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"wachiman/docker"
)

var (
	imageDangling bool
	imageForce    bool
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Gestiona las imágenes de Docker",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var imageLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lista todas las imágenes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		images, err := client.ListImages(imageDangling)
		if err != nil {
			return err
		}

		if len(images) == 0 {
			fmt.Println("No hay imágenes.")
			return nil
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, bold("IMAGE ID\tTAG\tSIZE\tCREATED\tCONTAINERS"))

		for _, img := range images {
			for i, tag := range img.Tags {
				id := ""
				if i == 0 {
					id = cyan(img.ID)
				}

				size := formatImageSize(img.Size)
				created := formatCreated(img.Created)

				containers := fmt.Sprintf("%d", img.Containers)
				if img.Containers == -1 {
					containers = "-"
				} else if img.Containers == 0 {
					containers = yellow("0")
				}

				tagStr := tag
				if tag == "<none>:<none>" {
					tagStr = red(tag)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id, tagStr, size, created, containers)
			}
		}
		w.Flush()
		return nil
	},
}

var imageRmCmd = &cobra.Command{
	Use:   "rm [imagen]",
	Short: "Elimina una imagen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		if err := client.RemoveImage(args[0], imageForce); err != nil {
			return err
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s imagen %s eliminada\n", green("✓"), args[0])
		return nil
	},
}

var imagePullCmd = &cobra.Command{
	Use:   "pull [imagen]",
	Short: "Pullea una imagen desde el registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := docker.New()
		if err != nil {
			return err
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("⬇ Pulleando %s...\n", cyan(args[0]))

		if err := client.PullImage(args[0]); err != nil {
			return err
		}

		fmt.Printf("%s %s descargada correctamente\n", green("✓"), args[0])
		return nil
	},
}

func formatImageSize(bytes int64) string {
	mb := float64(bytes) / 1024 / 1024
	if mb >= 1024 {
		return fmt.Sprintf("%.2f GB", mb/1024)
	}
	return fmt.Sprintf("%.0f MB", mb)
}

func formatCreated(unix int64) string {
	t := time.Unix(unix, 0)
	diff := time.Since(t)

	switch {
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	case diff < 365*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(diff.Hours()/24/30))
	default:
		return fmt.Sprintf("%d years ago", int(diff.Hours()/24/365))
	}
}

func init() {
	imageLsCmd.Flags().BoolVar(&imageDangling, "dangling", false, "Mostrar solo imágenes huérfanas")
	imageRmCmd.Flags().BoolVarP(&imageForce, "force", "f", false, "Forzar eliminación")

	ImageCmd.AddCommand(imageLsCmd)
	ImageCmd.AddCommand(imageRmCmd)
	ImageCmd.AddCommand(imagePullCmd)
}