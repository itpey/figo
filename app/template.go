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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type TemplateMetadata struct {
	Description string `json:"description"`
	Author      string `json:"author"`
}

func loadTemplateMetadata(templateDir string) (*TemplateMetadata, error) {
	metadataFile := filepath.Join(templateDir, metadataFileName)

	// Check if metadata file exists
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf(color.RedString("Error: checking metadata file: %v", err))
	}

	// Read metadata file
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf(color.RedString("Error: reading metadata file: %v", err))
	}

	// Unmarshal metadata JSON
	var metadata TemplateMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf(color.RedString("Error: decoding metadata file: %v", err))
	}

	return &metadata, nil
}
func extractAllTemplates(sourceDir string) error {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf(color.RedString("Error: reading repository directory: %v", err))
	}

	for _, file := range files {
		if file.IsDir() {
			templateDir := filepath.Join(sourceDir, file.Name())

			// Skip directories to be excluded
			if shouldSkipDir(file.Name()) {
				continue
			}

			// Read template metadata
			metaData, err := loadTemplateMetadata(templateDir)
			if err != nil {
				continue
			}

			// Skip if metadata is not available
			if metaData == nil {
				continue
			}

			// Check if the template already exists locally
			destPath := filepath.Join(templatesDirectory, file.Name())
			if _, err := os.Stat(destPath); err == nil {
				fmt.Print(color.YellowString("template '%s' already exists, updating to the latest version\n", file.Name()))
				if err := os.RemoveAll(destPath); err != nil {
					fmt.Print(color.RedString("Error: removing existing template '%s': %v\n", file.Name(), err))
					continue
				}
			}

			// Copy template directory to local templates directory
			if err := copyTemplate(templateDir, destPath, false); err != nil {
				fmt.Print(color.RedString("Error: copying template '%s': %v\n", file.Name(), err))
				continue
			}

			fmt.Print(color.GreenString("Template '%s' extracted successfully\n", file.Name()))
		}
	}

	return nil
}

func copyTemplate(src, dest string, excludeMetadata bool) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Traverse source directory
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Create directory in destination
			destPath := filepath.Join(dest, path[len(src):]) // Get relative path
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return err
			}
		} else {
			// Skip copying certain files
			if shouldSkipFile(info.Name()) {
				return nil
			}

			if excludeMetadata && info.Name() == metadataFileName {
				return nil
			}

			// Copy file to destination
			destPath := filepath.Join(dest, path[len(src):]) // Get relative path
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, srcFile); err != nil {
				return err
			}

			// Preserve file mode
			if err := os.Chmod(destPath, info.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
}

func printTemplateOptions(templates []string, currentIndex int) {
	for i, tmpl := range templates {
		templatePath := filepath.Join(templatesDirectory, tmpl)
		if metadata, err := loadTemplateMetadata(templatePath); err == nil {
			if i == currentIndex {
				fmt.Print(color.GreenString("-> %s: %s | %s\n", tmpl, metadata.Description, metadata.Author))
			} else {
				fmt.Printf("%s: %s | %s\n", tmpl, metadata.Description, metadata.Author)
			}

		}

	}
}

func downloadTemplates(url string) error {

	if url == "" {
		fmt.Print(color.RedString("No URL specified, using default repository: %s\n", defaultTemplatesRepoURL))
		url = defaultTemplatesRepoURL
	}

	fmt.Print(color.YellowString("downloading templates from repository: %s\n", url))

	repoDir := repoDirectory
	if err := gitClone(url, repoDir); err != nil {
		return fmt.Errorf(color.RedString("Error: cloning repository: %v", err))
	}
	defer os.RemoveAll(repoDir)

	if err := extractAllTemplates(repoDir); err != nil {
		return fmt.Errorf(color.RedString("Error: extracting templates: %v", err))
	}

	fmt.Println(color.GreenString("Templates downloaded successfully"))

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
				fmt.Println(color.YellowString("Select a template:"))
				printTemplateOptions(templates, currentIndex)
			}
		case keyboard.KeyArrowDown:
			if currentIndex < len(templates)-1 {
				currentIndex++
				clearConsole()
				fmt.Println(color.YellowString("Select a template:"))
				printTemplateOptions(templates, currentIndex)
			}
		case keyboard.KeyEnter:
			clearConsole()
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
