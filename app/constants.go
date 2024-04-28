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

import "github.com/urfave/cli/v2"

const (
	appNameArt = `______________               
___  ____/__(_)______ ______ 
__  /_   __  /__  __ /  __ \
_  __/   _  / _  /_/ // /_/ /
/_/      /_/  _\__, / \____/ 
              /____/         ` + "\n\n"
	appVersion     = "0.1.0"
	appDescription = `Figo is a command-line interface (CLI) tool designed to streamline the process of creating new Go projects. 
	With figo, you can quickly scaffold out the directory structure, configuration files, and necessary boilerplate code for your Go applications.`
	appCopyright            = "Apache-2.0 license\nFor more information, visit the GitHub repository: https://github.com/itpey/figo"
	defaultTemplatesRepoURL = "https://github.com/itpey/figo-templates"

	metaDataDirectoryname = "figo"
	metadataFileName      = "figo.json"
)

var (
	templatesDirectory = getDefaultDirectory("templates")
	repoDirectory      = getDefaultDirectory("repo")
)
var (
	appAuthors = []*cli.Author{{Name: "itpey", Email: "itpey@github.com"}}
)
