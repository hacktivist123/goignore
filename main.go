package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goignore",
		Short: "A lightweight CLI tool to generate .gitignore files",
	}

	rootCmd.AddCommand(newCmd)
	rootCmd.Execute()
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate and add .gitignore file to your project",
	Run: func(cmd *cobra.Command, args []string) {
		language, _ := cmd.Flags().GetString("language")

		if language == "" {
			fmt.Println("Please provide a programming language.")
			return
		}

		// Check if .git repo exists, if not initialize it
		_, err := os.Stat(".git")
		if err != nil {
			fmt.Println("Initializing a new Git repository...")
			err := execCommand("git", "init")
			if err != nil {
				fmt.Println("Error initializing Git repository:", err)
				return
			}
		}

		// Read .gitignore template content from file
		templateContent, err := readTemplateFile(language)
		if err != nil {
			fmt.Println("Error reading template file:", err)
			return
		}

		// Generate and write the .gitignore file
		err = generateGitignore(templateContent)
		if err != nil {
			fmt.Println("Error generating .gitignore:", err)
			return
		}

		fmt.Printf("Generated .gitignore for %s\n", language)
	},
}

func readTemplateFile(language string) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.txt", language)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func generateGitignore(content string) error {
	file, err := os.Create(".gitignore")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
