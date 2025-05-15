package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zkxjzmswkwl/workon/internal/wo"
)

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  add <projectName> <projectDir> <commandStr> [args...]")
	fmt.Println("    e.g `workon add projName E:\\Code\\Personal\\ProjPath code .`")
	fmt.Println("  config")
	fmt.Println("    e.g `workon config`")
	fmt.Println("  <projectName>")
	fmt.Println("    e.g `workon projName`")
	fmt.Println("  <projectName> cmd <cmdIndex> <shiftDir>")
	fmt.Println("    e.g `workon projName cmd 1 down` - moves command down, to the second index.")
	fmt.Println("    e.g `workon projName cmd 1 up` - moves command up, to the 0th index.")
	fmt.Println("  <projectName> cmd add <commandStr> [args...]")
	fmt.Println("    e.g `workon projName cmd add docker-compose up -d`")
	fmt.Println("  <projectName> cmd remove <cmdIndex>")
	fmt.Println("    e.g `workon projName cmd remove 1`")
}

func handleAddProject(config *wo.Config, args []string) {
	if len(args) < 5 {
		fmt.Println("Usage: workon add <projectName> <projectDir> <commandStr> [args...]")
		os.Exit(1)
	}
	projectName := args[2]
	projectDir := args[3]
	commandStr := args[4]
	config.CreateProject(projectName, projectDir, commandStr, args[5:]...)
}

func handleProjectDetails(config *wo.Config, projectName string) {
	project := config.GetProject(projectName)
	if project == nil {
		fmt.Println("Project not found")
		os.Exit(1)
	}
	fmt.Println("Details for", projectName)
	fmt.Println("  Project Name:", project.Name)
	fmt.Println("  Project Dir:", project.Dir)
	fmt.Println("Commands:")
	for i, command := range project.OpenFlow.Commands {
		fmt.Println("  ", i, command.Command, command.Args)
	}
	os.Exit(0)
}

func handleProjectCommandOperations(config *wo.Config, args []string, projectName string) {
	project := config.GetProject(projectName)
	if project == nil {
		fmt.Println("Project not found")
		os.Exit(1)
	}

	if len(args) <= 3 {
		fmt.Println("Usage: workon <projectName> cmd <operation> [args...]")
		os.Exit(1)
	}

	operation := args[3]

	if operation == "add" {
		if len(args) < 5 {
			fmt.Println("Usage: workon <projectName> cmd add <commandStr> [args...]")
			os.Exit(1)
		}
		commandStr := args[4]
		args := args[5:]
		config.InsertCommand(projectName, wo.CreateCommand(commandStr, args...))
		os.Exit(0)
	}

	if operation == "remove" {
		if len(args) < 5 {
			fmt.Println("Usage: workon <projectName> cmd remove <cmdIndex>")
			os.Exit(1)
		}
		cmdIndex, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println("Invalid command index")
			os.Exit(1)
		}
		config.RemoveCommand(projectName, cmdIndex)
		os.Exit(0)
	}

	cmdIndex, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("Invalid command index")
		os.Exit(1)
	}

	if len(args) < 5 {
		fmt.Println("Usage: workon <projectName> cmd <cmdIndex> <up|down>")
		os.Exit(1)
	}

	shiftDir := args[4]
	if shiftDir == "down" && cmdIndex < len(project.OpenFlow.Commands)-1 {
		project.OpenFlow.Commands[cmdIndex], project.OpenFlow.Commands[cmdIndex+1] =
			project.OpenFlow.Commands[cmdIndex+1], project.OpenFlow.Commands[cmdIndex]
		config.Save()
		fmt.Println("Command moved down")
	} else if shiftDir == "up" && cmdIndex > 0 {
		project.OpenFlow.Commands[cmdIndex], project.OpenFlow.Commands[cmdIndex-1] =
			project.OpenFlow.Commands[cmdIndex-1], project.OpenFlow.Commands[cmdIndex]
		config.Save()
		fmt.Println("Command moved up")
	} else {
		fmt.Println("Cannot move command (invalid direction or at boundary)")
	}
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: workon <command>")
		fmt.Println("  e.g `workon help`")
		os.Exit(1)
	}

	config := wo.GetOrCreateConfig()
	firstArg := os.Args[1]

	switch firstArg {
	case "help":
		printHelp()
		os.Exit(0)
	case "add":
		handleAddProject(config, os.Args)
	case "config":
		config.Print()
		os.Exit(0)
	default:
		projectName := firstArg

		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "details":
				handleProjectDetails(config, projectName)
			case "cmd":
				handleProjectCommandOperations(config, os.Args, projectName)
			}
		}

		project := config.GetProject(projectName)
		if project == nil {
			fmt.Println("Project not found")
			os.Exit(1)
		}
		fmt.Println("Opening project", projectName)
		project.Open()
	}
}
