//go:generate go run github.com/hacktivist123/goignore/
package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var extensions = map[string][]string{
	"golang":     {".go"},
	"javascript": {".js", ".ts", ".tsx"},
	"python":     {".py"},
	"c++":        {".cpp", ".h"},
	"rust":       {".rs"},
	"ruby":       {".rb"},
	"c":          {".c"},
	"haskell":    {".hs"},
}

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
		if language == "" || autoDetect {
			language = detectLanguage()
			if language == "" {
				color.Red("Unable to auto-detect programming language. Please provide a language manually.")
				return
			}
		}

		if language == "" {
			color.Red("Please provide a programming language.")
			return
		}

		// // Check if .git repo exists, if not initialize it
		// _, err := os.Stat(".git")
		// if err != nil {
		// 	color.Yellow("Initializing a new Git repository...")
		// 	err := execCommand("git", "init")
		// 	if err != nil {
		// 		color.Red("Error initializing Git repository:", err)
		// 		return
		// 	}
		// }

		// Read .gitignore template content from file
		templateContent, err := readTemplateFile(language)

		// if the language template does not exist
		if errors.Is(err, fs.ErrNotExist) {
			err = fmt.Errorf("language '%s' not supported", language)
		}

		if err != nil {
			color.Red("Error: %s", err)
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
	// struct to store file extension
	languagePercentage := make(map[string]int)

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		for lang, exts := range extensions {
			for _, e := range exts {
				if ext == e {
					if languagePercentage[lang] > 0 {
						languagePercentage[lang]++
					} else {
						languagePercentage[lang] = 1
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return ""
	}
	// get the one with highest occurrence
	fileExt := highestOccurrence(languagePercentage)
	return fileExt
}

func highestOccurrence(data map[string]int) string {
	language := ""
	languageValue := 0
	for key, value := range data {
		if value > languageValue {
			languageValue = value
			language = key
		}
	}
	return language
}

func getSupportedLanguages() []string {
	result := []string{}

	for lang := range extensions {
		result = append(result, lang)
	}
	return result
}
func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)

	// Define the 'language' and 'auto-detect' flags for the newCmd
	newCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language for .gitignore file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//go:embed templates/*.txt
var templateFiles embed.FS

func readTemplateFile(language string) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.txt", language)
	content, err := templateFiles.ReadFile(templatePath)
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

// func execCommand(command string, args ...string) error {
// 	cmd := exec.Command(command, args...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }
