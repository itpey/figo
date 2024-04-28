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

func initializeGitRepository(projectPath string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	return cmd.Run()
}

func runGoGet(projectPath string) error {
	cmd := exec.Command("go", "get", "./...")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(color.RedString("Error: go get failed: %v\n%s", err, output))
	}
	return nil
}

func runGoModTidy(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(color.RedString("Error: go mod tidy failed: %v\n%s", err, output))
	}
	return nil
}

func gitClone(repoURL, destination string) error {
	cmd := exec.Command("git", "clone", repoURL, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(color.RedString("Error: git clone error: %v, output: %s", err, string(output)))
	}
	return nil
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
