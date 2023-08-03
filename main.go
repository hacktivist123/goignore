package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var language string
var autoDetect bool

var rootCmd = &cobra.Command{
	Use:   "goignore",
	Short: "A lightweight CLI tool for generating .gitignore files",
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate and add .gitignore file to your project",
	Run: func(cmd *cobra.Command, args []string) {
		if language == "" && !autoDetect {
			color.Red("Please provide a programming language.")
			return
		}

		if autoDetect {
			language = detectLanguage()
			if language == "" {
				color.Red("Unable to auto-detect programming language. Please provide a language manually.")
				return
			}
		}

		// Check if .git repo exists, if not initialize it
		_, err := os.Stat(".git")
		if err != nil {
			color.Yellow("Initializing a new Git repository...")
			err := execCommand("git", "init")
			if err != nil {
				color.Red("Error initializing Git repository:", err)
				return
			}
		}

		// Read .gitignore template content from file
		templateContent, err := readTemplateFile(language)
		if err != nil {
			color.Red("Error reading template file:", err)
			return
		}

		// Generate and write the .gitignore file
		err = generateGitignore(templateContent)

		if err != nil {
			color.Red("Error generating .gitignore:", err)
			return
		}

		path, err := filepath.Abs(".gitignore")
		if err != nil {
			color.Red("Error getting absolute path:", err)
			return
		}

		color.Green("Generated .gitignore for %s in path: %s", language, path)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported programming languages",
	Run: func(cmd *cobra.Command, args []string) {
		supportedLanguages := getSupportedLanguages()
		color.Cyan("Supported programming languages:")
		for _, lang := range supportedLanguages {
			fmt.Println("-", lang)
		}
	},
}

func detectLanguage() string {
	files, err := os.ReadDir(".")
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".go" {
				return "golang"
			} else if ext == ".js" || ext == ".ts" || ext == ".tsx" {
				return "javaScript"
			} else if ext == ".py" {
				return "python"
			} else if ext == ".cpp" || ext == ".h" {
				return "c++"
			}
		}
	}

	return ""
}

func getSupportedLanguages() []string {
	// Add any additional supported languages here
	return []string{"python", "javascript", "golang", "c++"}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)

	// Define the 'language' flag
	newCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language for .gitignore file (use 'auto' for auto-detection)")

	// Set PersistentPreRun function to handle flag validation and auto-detection
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if language == "" {
			fmt.Println("Error: Please provide a programming language or use 'auto' for auto-detection.")
			os.Exit(1)
		}

		if language == "auto" {
			language = detectLanguage()
			if language == "" {
				fmt.Println("Error: Unable to auto-detect programming language. Please provide a language manually.")
				os.Exit(1)
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
