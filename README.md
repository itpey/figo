<div align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/itpey/figo/main/static/images/figo_icon_dark.png"  width="250" height="100">
  <img alt="figo" src="https://raw.githubusercontent.com/itpey/figo/main/static/images/figo_icon-light.png"  width="250" height="100" >
</picture>
</div>

<p align="center">
Figo is a command-line tool for rapidly scaffolding new Go projects based on customizable templates. It provides various commands to create, manage, and work with project templates efficiently.
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/itpey/keycontrol">
    <img src="https://pkg.go.dev/badge/github.com/itpey/figo.svg" alt="Go Reference">
  </a>
  <a href="https://github.com/itpey/figo/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/itpey/figo" alt="license">
  </a>
</p>

# Features

- **Project Creation**: Create new Go projects from predefined templates.
- **Template Management**:
  - List available project templates.
  - Add templates from Git repositories.
  - Delete specific or all project templates.
- **Environment Check**: Verify system environment for required tools (Git and Go).

# Installation

Make sure you have Go installed and configured on your system. Use go install to install Figo:

```bash
go install github.com/itpey/figo@latest
```

Ensure that your `GOBIN` directory is in your `PATH` for the installed binary to be accessible globally.

# Usage

## Checking System Environment

To check the system environment for required tools (Git and Go):

```bash
figo doctor
```

## Creating a New Project

To create a new Go project with Figo:

```bash
figo
```

This will prompt you to enter the project name and select a template interactively.

## Advanced Usage

To view detailed usage instructions and available commands:

```bash
figo -h
```

This will display comprehensive information about using Figo, including commands, options, and examples.

# Feedback and Contributions

If you encounter any issues or have suggestions for improvement, please [open an issue](https://github.com/itpey/figo/issues) on GitHub.

We welcome contributions! Fork the repository, make your changes, and submit a pull request.

# License

Figo is open-source software released under the Apache License, Version 2.0. You can find a copy of the license in the [LICENSE](https://github.com/itpey/figo/blob/main/LICENSE) file.

# Author

Figo was created by [itpey](https://github.com/itpey).
