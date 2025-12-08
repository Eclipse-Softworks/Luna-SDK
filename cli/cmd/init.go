package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Luna project",
	Long:  `Scaffold a new project with Luna SDK integration. Supports TypeScript, Python, and Go.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "my-luna-app"
		if len(args) > 0 {
			projectName = args[0]
		}

		lang, _ := cmd.Flags().GetString("lang")
		validLangs := map[string]bool{"ts": true, "python": true, "go": true}
		if !validLangs[lang] {
			return fmt.Errorf("invalid language '%s'. Supported: ts, python, go", lang)
		}

		fmt.Printf("Initiailzing new %s project '%s'...\n", lang, projectName)

		if err := os.Mkdir(projectName, 0755); err != nil {
			if !os.IsExist(err) {
				return err
			}
		}

		var err error
		switch lang {
		case "ts":
			err = scaffoldTypeScript(projectName)
		case "python":
			err = scaffoldPython(projectName)
		case "go":
			err = scaffoldGo(projectName)
		}

		if err != nil {
			return fmt.Errorf("failed to scaffold project: %w", err)
		}

		fmt.Printf("\n✓ Project created in ./%s\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)
		switch lang {
		case "ts":
			fmt.Println("  npm install")
			fmt.Println("  npm start")
		case "python":
			fmt.Println("  pip install -r requirements.txt")
			fmt.Println("  python main.py")
		case "go":
			fmt.Println("  go mod tidy")
			fmt.Println("  go run main.go")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("lang", "l", "ts", "Language to use (ts, python, go)")
}

func scaffoldTypeScript(name string) error {
	packageJson := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "main": "index.js",
  "dependencies": {
    "@eclipse/luna-sdk": "^1.0.0",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/node": "^18.0.0",
    "ts-node": "^10.0.0"
  },
  "scripts": {
    "start": "ts-node src/index.ts"
  }
}`, name)

	indexTs := `import { LunaClient } from '@eclipse/luna-sdk';
import * as dotenv from 'dotenv';

dotenv.config();

const client = new LunaClient({
  apiKey: process.env.LUNA_API_KEY
});

async function main() {
  try {
    const users = await client.users.list({ limit: 5 });
    console.log('Users:', users.data);
  } catch (err) {
    console.error('Error:', err);
  }
}

main();`

	if err := writeFile(filepath.Join(name, "package.json"), packageJson); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(name, "src"), 0755); err != nil {
		return err
	}
	return writeFile(filepath.Join(name, "src", "index.ts"), indexTs)
}

func scaffoldPython(name string) error {
	requirements := `luna-sdk>=1.0.0
python-dotenv>=1.0.0`

	mainPy := `import os
import asyncio
from dotenv import load_dotenv
from luna import LunaClient

load_dotenv()

async def main():
    client = LunaClient(api_key=os.getenv("LUNA_API_KEY"))
    
    try:
        users = await client.users.list(limit=5)
        print(f"Users: {users.data}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    asyncio.run(main())`

	if err := writeFile(filepath.Join(name, "requirements.txt"), requirements); err != nil {
		return err
	}
	return writeFile(filepath.Join(name, "main.py"), mainPy)
}

func scaffoldGo(name string) error {
	modName := strings.ReplaceAll(strings.ToLower(name), " ", "-")

	mainGo := `package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
)

func main() {
	apiKey := os.Getenv("LUNA_API_KEY")
	if apiKey == "" {
		log.Fatal("LUNA_API_KEY not set")
	}

	client := luna.NewClient(luna.WithAPIKey(apiKey))

	users, err := client.Users().List(context.Background(), &luna.ListParams{Limit: 5})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Users: %+v\n", users.Data)
}`

	if err := runGoModInit(name, modName); err != nil {
		return err
	}

	return writeFile(filepath.Join(name, "main.go"), mainGo)
}

func runGoModInit(dir, modName string) error {
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
