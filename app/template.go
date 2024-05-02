// Copyright 2024 itpey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

func extractAllTemplates(sourceDir string, repoName string) error {

	if isGoModule(sourceDir) {
		destPath := filepath.Join(templatesDirectory, repoName)
		if err := copyTemplate(sourceDir, destPath); err != nil {
			fmt.Printf(color.RedString("Error:Error copying template '%s': %v\n"), repoName, err)
		}
		fmt.Printf(color.GreenString("Template '%s' extracted successfully\n"), repoName)
	}

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf(color.RedString("Error: reading repository directory: %v"), err)
	}

	for _, file := range files {
		if file.IsDir() {
			templateDir := filepath.Join(sourceDir, file.Name())

			// Skip directories to be excluded
			if shouldSkipDir(file.Name()) {
				continue
			}

			// Check if the directory contains a Go module
			if isGoModule(templateDir) {
				// Get the base directory name (to use as template name)
				templateName := fmt.Sprintf("%s_%s", repoName, file.Name())

				// Copy template directory to local templates directory
				destPath := filepath.Join(templatesDirectory, templateName)
				if err := copyTemplate(templateDir, destPath); err != nil {
					fmt.Printf(color.RedString("Error: copying template '%s': %v\n"), templateName, err)
					continue
				}

				fmt.Printf(color.GreenString("Template '%s' extracted successfully\n"), templateName)
			}
		}
	}

	return nil
}
func copyTemplate(src, dest string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf(color.RedString("Error: failed to create destination directory: %v"), err)
	}

	// Traverse source directory
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(color.RedString("Error: accessing path %q: %v"), path, err)
		}

		// Get relative path within the source directory
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf(color.RedString("Error: getting relative path for %q: %v"), path, err)
		}

		// Construct destination path
		destPath := filepath.Join(dest, relPath)
		// Check if the directory should be skipped
		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			// Create directory in destination
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf(color.RedString("Error: failed to create directory %q: %v"), destPath, err)
			}
		} else {
			if shouldSkipFile(info.Name()) {
				return nil
			}
			// Copy file to destination
			srcFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf(color.RedString("Error: failed to open source file %q: %v"), path, err)
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf(color.RedString("Error: failed to create destination file %q: %v"), destPath, err)
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, srcFile); err != nil {
				return fmt.Errorf(color.RedString("Error: failed to copy file %q to %q: %v"), path, destPath, err)
			}

			// Preserve file mode
			if err := os.Chmod(destPath, info.Mode()); err != nil {
				return fmt.Errorf(color.RedString("Error: failed to set file mode for %q: %v"), destPath, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf(color.RedString("Error: copying template: %v"), err)
	}

	return nil
}

func printTemplateOptions(templates []string, currentIndex int) {
	for i, tmpl := range templates {
		if i == currentIndex {
			fmt.Print(color.GreenString("-> %s\n", tmpl))
		} else {
			fmt.Printf("%s\n", tmpl)
		}

	}

}

func downloadTemplates(url string) error {

	if url == "" {
		fmt.Print(color.YellowString("Warning: No URL specified, using default repository: %s\n", defaultTemplatesRepoURL))
		url = defaultTemplatesRepoURL
	}

	fmt.Print(color.YellowString("Downloading templates from repository: %s ...\n", url))

	repoName, err := extractRepoNameFromURL(url)
	if err != nil {
		return fmt.Errorf(color.RedString("Error: extracting repository name: %v"), err)
	}

	repoDir := repoDirectory
	if err := gitClone(url, repoDir); err != nil {
		return fmt.Errorf(color.RedString("Error: cloning repository: %v", err))
	}
	defer os.RemoveAll(repoDir)

	if err := extractAllTemplates(repoDir, repoName); err != nil {
		return fmt.Errorf(color.RedString("Error: extracting templates: %v", err))
	}

	return nil
}

func deleteAllTemplates() error {

	if _, err := os.Stat(templatesDirectory); os.IsNotExist(err) {
		return fmt.Errorf(color.RedString("Error: templates directory does not exist"))
	}

	err := os.RemoveAll(templatesDirectory)
	if err != nil {
		return fmt.Errorf(color.RedString("Error: deleting templates directory: %v", err))
	}

	fmt.Println(color.GreenString("All templates deleted successfully"))
	return nil
}

func deleteTemplateByName(templateName string) error {

	templatePath := filepath.Join(templatesDirectory, templateName)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf(color.RedString("Error: template '%s' not found", templateName))
	}

	err := os.RemoveAll(templatePath)
	if err != nil {
		return fmt.Errorf(color.RedString("Error: deleting template: %v", err))

	}
	fmt.Print(color.GreenString("template '%s' deleted successfully\n", templateName))
	return nil
}

func listTemplates() ([]string, error) {

	if _, err := os.Stat(templatesDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(templatesDirectory, 0755); err != nil {
			return nil, fmt.Errorf(color.RedString("Error: creating templates directory: %v", err))
		}
		if err := downloadTemplates(defaultTemplatesRepoURL); err != nil {
			return nil, err
		}

	}

	var templates []string

	files, err := os.ReadDir(templatesDirectory)
	if err != nil {
		return nil, fmt.Errorf(color.RedString("Error: reading templates directory: %v", err))
	}

	for _, file := range files {
		if file.IsDir() {
			templates = append(templates, file.Name())
		}
	}

	return templates, nil
}

func selectTemplate(templates []string) (string, error) {
	clearConsole()
	fmt.Print(color.CyanString(appNameArt))
	fmt.Println(color.YellowString("Select a template:"))

	currentIndex := 0
	printTemplateOptions(templates, currentIndex)

	err := keyboard.Open()
	if err != nil {
		return "", fmt.Errorf(color.RedString("Error: opening keyboard: %v", err))

	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return "", fmt.Errorf(color.RedString("Error: reading keyboard input: %v", err))
		}

		switch key {
		case keyboard.KeyArrowUp:
			if currentIndex > 0 {
				currentIndex--
				clearConsole()
				fmt.Print(color.CyanString(appNameArt))

				fmt.Println(color.YellowString("Select a template:"))
				printTemplateOptions(templates, currentIndex)
			}
		case keyboard.KeyArrowDown:
			if currentIndex < len(templates)-1 {
				currentIndex++
				clearConsole()
				fmt.Print(color.CyanString(appNameArt))
				fmt.Println(color.YellowString("Select a template:"))
				printTemplateOptions(templates, currentIndex)
			}
		case keyboard.KeyEnter:
			clearConsole()
			fmt.Print(color.CyanString(appNameArt))
			selectedTemplate := templates[currentIndex]
			return selectedTemplate, nil
		case keyboard.KeyEsc:
			return "", fmt.Errorf(color.RedString("Error: selection canceled."))
		}

		if char == rune(keyboard.KeyEsc) {
			break
		}
	}

	return "", nil
}
