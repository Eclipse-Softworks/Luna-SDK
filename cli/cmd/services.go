package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage services",
	Long:  `List and manage available specific services within the Eclipse Softworks Platform.`,
}

var servicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tSTATUS")
		fmt.Fprintln(w, "----\t-----------\t------")

		services := []struct {
			Name, Desc, Status string
		}{
			{"Identity", "User and group management", "Active"},
			{"ResMate", "Student residence listings", "Active"},
			{"Storage", "Cloud object storage", "Active"},
			{"AI Objects", "Generative AI completions", "Beta"},
			{"Automation", "Workflow automation", "Preview"},
		}

		for _, s := range services {
			fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.Desc, s.Status)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	servicesCmd.AddCommand(servicesListCmd)
}
