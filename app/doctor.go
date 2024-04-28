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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func runDoctor() error {

	// Check Git version
	gitVersionCmd := exec.Command("git", "--version")
	gitVersionOutput, err := gitVersionCmd.CombinedOutput()
	if err != nil {
		fmt.Print(color.RedString("Error: git is not installed or not available in PATH"))
	} else {
		fmt.Print(color.GreenString(string(gitVersionOutput)))
	}

	// Check Go version
	goVersionCmd := exec.Command("go", "version")
	goVersionOutput, err := goVersionCmd.CombinedOutput()
	if err != nil {
		fmt.Print(color.RedString("Error: go is not installed or not available in PATH"))
	} else {
		fmt.Print(color.GreenString(string(goVersionOutput)))
	}

	// Check if both Git and Go are available
	if err == nil {
		fmt.Println(color.GreenString("All required tools are installed and accessible."))
		// Check if template directory is empty
		templates, err := listTemplates()
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			fmt.Println(color.YellowString("Warning: No templates found in the templates directory."))

			// Prompt user to download default templates
			fmt.Println("Do you want to download default templates? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "y" || answer == "Y" {
				err := downloadTemplates(defaultTemplatesRepoURL)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(color.YellowString("No templates downloaded."))
			}
		}
		fmt.Println(color.YellowString("Run 'figo help' for usage instructions."))
	} else {
		fmt.Println(color.RedString("Error: system environment check completed with errors."))
		fmt.Println(color.YellowString("Please ensure that Git and Go are installed and available in your PATH."))
	}

	return nil
}
