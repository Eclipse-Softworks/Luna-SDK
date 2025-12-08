package cmd

import (
	"fmt"
	"net/http"
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

		// Connectivity Check
		fmt.Printf("\nConnectivity:\n")
		// Using a public endpoint or the base URL from constants if possible.
		// We'll hardcode the known production generic endpoint for health check.
		resp, err := http.Get("https://api.eclipse.dev/health")
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("  ✓ API Gateway (api.eclipse.dev): Reachable")
		} else {
			fmt.Println("  ✗ API Gateway (api.eclipse.dev): Unreachable or Down")
		}

		// check Auth server
		resp, err = http.Get("https://auth.eclipse.dev/.well-known/openid-configuration")
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("  ✓ Auth Service (auth.eclipse.dev): Reachable")
		} else {
			fmt.Println("  ✗ Auth Service (auth.eclipse.dev): Unreachable")
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
