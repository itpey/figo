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
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func showVersion() error {
	fmt.Printf("figo version %s\n", appVersion)
	return nil
}

func getDefaultDirectory(name string) string {
	var templatesDir string

	switch runtime.GOOS {
	case "windows":
		// For Windows, use %APPDATA%\figo\...
		templatesDir = filepath.Join(os.Getenv("APPDATA"), metaDataDirectoryname, name)
	default:
		// For Unix-like systems (Linux, macOS), use ~/.config/figo/...
		templatesDir = filepath.Join(os.Getenv("HOME"), ".config", metaDataDirectoryname, name)
	}

	return templatesDir
}

func runCommand(command string, args []string, projectPath string, description string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = projectPath
	fmt.Printf("Running %s ...\n", description)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(color.RedString("Error: running %s: %v\n%s", description, err, output))
	}

	fmt.Printf(color.GreenString("%s finished successfully.\n"), description)
	return nil
}

// Function to initialize a Git repository in the specified project path
func initializeGitRepository(projectPath string) error {
	return runCommand("git", []string{"init"}, projectPath, "git init")
}

// Function to run 'go get ./...' in the specified project path
func runGoGet(projectPath string) error {
	return runCommand("go", []string{"get"}, projectPath, "go get")
}

// Function to run 'go mod tidy' in the specified project path
func runGoModTidy(projectPath string) error {
	return runCommand("go", []string{"mod", "tidy"}, projectPath, "go mod tidy")
}

// Function to clone a Git repository from the specified URL to the destination directory
func gitClone(repoURL, destination string) error {
	return runCommand("git", []string{"clone", repoURL, destination}, ".", "git clone")
}

func promptProjectName() string {
	fmt.Print(color.YellowString("Input your project name: "))
	var projectName string
	fmt.Scanln(&projectName)
	return strings.TrimSpace(projectName)
}

func shouldSkipFile(fileName string) bool {
	ignoredFiles := map[string]bool{}
	return ignoredFiles[fileName]
}

func shouldSkipDir(dirName string) bool {
	skipDirs := map[string]bool{
		".git":    true,
		".github": true,
	}
	return skipDirs[dirName]
}

func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}

	for _, char := range name {
		if !isValidCharacter(char) {
			return false
		}
	}

	return true
}
func isValidCharacter(char rune) bool {

	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-' ||
		char == '_'
}

func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func isGoModule(dir string) bool {
	goModFile := filepath.Join(dir, "go.mod")

	// Use os.Stat to check the existence of the go.mod file
	if _, err := os.Stat(goModFile); err != nil {
		if os.IsNotExist(err) {
			return false // go.mod file does not exist
		}
		// Handle other errors (e.g., permission denied)
		return false
	}

	return true // go.mod file exists
}

func extractRepoNameFromURL(repoURL string) (string, error) {
	// Parse the repository URL
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf(color.RedString("Error: failed to parse repository URL: %v"), err)
	}

	// Remove .git extension and split the path into segments
	pathSegments := strings.Split(strings.TrimSuffix(parsedURL.Path, ".git"), "/")

	// Find the last non-empty segment in the path
	var repoName string
	for i := len(pathSegments) - 1; i >= 0; i-- {
		if pathSegments[i] != "" {
			repoName = pathSegments[i]
			break
		}
	}

	// Validate and return the repository name
	if repoName == "" {
		return "", fmt.Errorf(color.RedString("Error: unable to determine repository name from URL: %s"), repoURL)
	}
	return repoName, nil
}
