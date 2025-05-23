# workon

A simple CLI tool for managing development projects and automating their startup processes.


## Installation

```
go install github.com/zkxjzmswkwl/workon/cmd/workon
```

## Config location
Config is located at `~/.wo/wo_config.json`

## Usage

```
workon <command>
```

### Commands

- `add <projectName> <projectDir> <commandStr> [args...]`  
  Add a new project with an initial command
  
- `config`  
  Display current configuration

- `<projectName>`  
  Open the specified project and execute its commands

- `<projectName> details`  
  Show details about a specific project

- `<projectName> cmd <cmdIndex> <shiftDir>`  
  Reorder commands for a project (move up or down)

- `<projectName> cmd add <commandStr> [args...]`  
  Add command

- `<projectName> cmd remove <cmdIndex>`</br>
  Remove command.

- `help`  
  Display help information

## Examples

```
# Create a new project
workon add myapp E:\Code\Projects\myapp code .

# Add a command
workon myapp cmd add docker-compose up -d

# Open project
workon myapp

# View project details
workon myapp details

# Reorder commands
workon myapp cmd 1 up
``` 

## Config example
```json
{
  "Projects": [
    {
      "Timestamp": 1747268259,
      "Name": "pygzr",
      "Dir": "E:\\py-gzr",
      "OpenFlow": {
        "Timestamp": 1747268259,
        "Commands": [
          {
            "Command": "code",
            "Args": [
              "."
            ]
          }
        ]
      }
    },
    {
      "Timestamp": 1747269978,
      "Name": "igz",
      "Dir": "C:\\Users\\dev\\proj\\src",
      "OpenFlow": {
        "Timestamp": 1747269978,
        "Commands": [
          {
            "Command": "cmd",
            "Args": [
              "/C",
              "start",
              "",
              "C:\\Program Files (x86)\\International GunZ\\client\\"
            ]
          },
          {
            "Command": "cmd",
            "Args": [
              "/C",
              "start",
              "",
              "C:\\Users\\Dev\\proj\\src\\build\\bin\\RelWithDebInfo"
            ]
          },
          {
            "Command": "clion64",
            "Args": [
              "."
            ]
          }
        ]
      }
    }
  ]
}
```