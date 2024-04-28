package app

import (
	"fmt"
	"os/exec"

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
		fmt.Print(color.YellowString("Run 'figo help' for usage instructions."))
	} else {
		fmt.Println(color.RedString("Error: system environment check completed with errors."))
		fmt.Print(color.YellowString("Please ensure that Git and Go are installed and available in your PATH."))
	}

	return nil
}
