package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check SDK environment",
	Long:  "Diagnose your environment and ensure requirements for Luna SDK are met.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Luna SDK Diagnostics")
		fmt.Println("====================")

		// System Info
		fmt.Printf("\nSystem:\n")
		fmt.Printf("  OS: %s\n", runtime.GOOS)
		fmt.Printf("  Arch: %s\n", runtime.GOARCH)

		// Language Checks
		fmt.Printf("\nEnvironment:\n")
		checkCommand("go", "version")
		checkCommand("node", "--version")
		checkCommand("npm", "--version")
		checkCommand("python", "--version")

		// Auth Check
		fmt.Printf("\nAuthentication:\n")
		apiKey := getAPIKey()
		if apiKey != "" {
			fmt.Printf("  ✓ API Key found (%s...)\n", apiKey[:5])
		} else {
			fmt.Println("  ✗ No API Key found (Run 'luna auth login')")
		}

		fmt.Println("\nDiagnostics complete.")
	},
}

func checkCommand(name string, arg string) {
	cmd := exec.Command(name, arg)
	out, err := cmd.Output()
	if err == nil {
		fmt.Printf("  ✓ %s: %s", name, out)
	} else {
		fmt.Printf("  ✗ %s: Not found\n", name)
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
