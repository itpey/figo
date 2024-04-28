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
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func Create() *cli.App {
	app := &cli.App{
		Name:        "figo",
		Usage:       "A CLI tool for rapidly scaffolding new Go projects",
		Description: appDescription,
		Copyright:   appCopyright,
		Authors:     appAuthors,
		Version:     appVersion,
		Action: func(c *cli.Context) error {
			clearConsole()
			fmt.Print(color.CyanString(appNameArt))

			var projectName string

			for {
				projectName = promptProjectName()
				if projectName == "" {
					fmt.Println(color.RedString("Error: project name cannot be empty"))
					continue
				}

				if !isValidProjectName(projectName) {
					fmt.Println(color.RedString("Error: invalid project name: %s\n", projectName))
					continue
				}
				break
			}

			templates, err := listTemplates()
			if err != nil {
				return err
			}

			if len(templates) == 0 {
				return fmt.Errorf(color.RedString("Error: no templates found in", templatesDirectory))
			}

			selectedTemplate, err := selectTemplate(templates)
			if err != nil {
				return err
			}

			createProject(projectName, selectedTemplate)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "doctor",
				Usage:   "Check system environment for Git and Go",
				Aliases: []string{"doc"},
				Action: func(c *cli.Context) error {
					return runDoctor()
				},
			},
			{
				Name:    "create",
				Aliases: []string{"init", "new", "i", "c"},
				Usage:   "Create a new Go project",
				Action: func(c *cli.Context) error {
					return createProject(c.String("name"), c.String("template"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name of the project",
						Required: true,
					},

					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Usage:   "Project template to use",
						Value:   "default",
					},
				},
			},
			{
				Name:    "list-templates",
				Aliases: []string{"lt", "ls", "l"},
				Usage:   "List available figo project templates",
				Action: func(c *cli.Context) error {

					templates, err := listTemplates()
					if err != nil {
						return err
					}

					if len(templates) == 0 {
						return fmt.Errorf(color.RedString("Error: no templates found in", templatesDirectory))
					}

					fmt.Println(color.YellowString("Available templates:"))
					printTemplateOptions(templates, -1)
					return nil

				},
			},
			{
				Name:    "add-templates",
				Aliases: []string{"add"},
				Usage:   "Download figo project templates from a Git repository",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u"},
						Usage:   "Git repository URL to download templates from",
					},
				},
				Action: func(c *cli.Context) error {
					url := c.String("url")
					return downloadTemplates(url)

				},
			},
			{
				Name:    "delete-templates",
				Aliases: []string{"del"},
				Usage:   "Delete figo project templates",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Required: true,
						Usage:    "Delete a specific figo project template by name",
						Action: func(c *cli.Context, name string) error {
							return deleteTemplateByName(name)
						},
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "all",
						Usage: "Delete all figo project templates",
						Action: func(c *cli.Context) error {
							return deleteAllTemplates()
						},
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v", "ver", "about"},
				Usage:   "Print the version",
				Action: func(c *cli.Context) error {
					return showVersion()
				},
			},
		},
	}

	return app
}

func createProject(projectName string, templateName string) error {
	fmt.Print(color.YellowString("Creating project '%s'...\n", projectName))

	templatePath := filepath.Join(templatesDirectory, templateName)

	if _, err := os.Stat(templatePath); err != nil {
		return fmt.Errorf(color.RedString("Error: template '%s' not found", templateName))
	}

	projectPath := filepath.Join(".", projectName)

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf(color.RedString("Error: creating project directory: %v", err))
	}

	// Copy project from template directory to project directory
	if err := copyTemplate(templatePath, projectPath); err != nil {
		return fmt.Errorf(color.RedString("Error: copying project: %v", err))
	}

	// Initialize Git repository
	if err := initializeGitRepository(projectPath); err != nil {
		return fmt.Errorf(color.RedString("Error: initializing Git repository: %v", err))
	}

	// Run 'go get' to fetch dependencies (if any)
	if err := runGoGet(projectPath); err != nil {
		return fmt.Errorf(color.RedString("Error: running 'go get': %v", err))
	}

	// Run 'go mod tidy' to clean up go.mod and go.sum files
	if err := runGoModTidy(projectPath); err != nil {
		return fmt.Errorf(color.RedString("Error: running 'go mod tidy': %v", err))
	}

	fmt.Print(color.GreenString("Project '%s' created successfully!\n", projectName))

	return nil
}
