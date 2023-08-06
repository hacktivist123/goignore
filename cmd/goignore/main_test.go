package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestDetectLanguage(t *testing.T) {
	testFs := afero.NewMemMapFs()
	if err := testFs.Mkdir("src/", 0755); err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	file, err := testFs.Create("src/test.go")
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	
	t.Log(testFs.Name())

	want := "golang"
	answer := detectLanguage(testFs)

	if want != answer {
		t.Fatalf("Wanted %s, got %s", want, answer)
	}
}

func TestGenerateGitignore(t *testing.T) {
	testFs := afero.NewMemMapFs()
	content := ".env"
	err := generateGitignore(testFs, content)
	if err != nil {
		panic(err)
	}
	want := true
	exists, _ := afero.Exists(testFs, ".gitignore")
	if !exists {
		t.Fatalf("Wanted %v, got %v", want, exists)
	}

}

func TestGetSupportedLanguages(t *testing.T) {
	languages := getSupportedLanguages()
	want := len(languages)
	answer := len(extensions)
	if want != answer {
		t.Fatalf("Wanted %d, got %d", want, answer)
	}
}

func TestReadTemplateFile(t *testing.T) {
	content, _ := templateFiles.ReadFile("templates/golang.txt")
	want := string(content)

	answer, err := readTemplateFile("golang")
	if err != nil {
		panic(err)
	}
	if want != answer {
		t.Fatalf("Wanted %s, got %s", want, answer)

	}

	// Test for python
	content, _ = templateFiles.ReadFile("templates/python.txt")
	want = string(content)

	answer, err = readTemplateFile("python")
	if err != nil {
		panic(err)
	}
	if want != answer {
		t.Fatalf("Wanted %s, got %s", want, answer)

	}

}
