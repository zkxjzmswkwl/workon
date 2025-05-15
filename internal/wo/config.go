package wo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Helper - Run command from given working directory
func RunCommand(command string, dir string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	fmt.Println("Command", command, "run in", dir)
	fmt.Println("Output:", string(output))
}

func GetOrCreateConfig() *Config {
	if DoesConfigExist() {
		return LoadConfig()
	}
	return NewConfig()
}

func ensureConfigFolderExists() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configFolder := filepath.Join(homeDir, ".wo")
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		os.Mkdir(configFolder, 0755)
	}
}

// Checks if the config both exists on the filesystem
// and contains at least one byte.
func DoesConfigExist() bool {
	ensureConfigFolderExists()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	configFolder := filepath.Join(homeDir, ".wo")

	configPath := filepath.Join(configFolder, "wo_config.json")
	fInfo, err := os.Stat(configPath)
	if err != nil {
		return false
	}
	if fInfo.Size() == 0 {
		fmt.Println("Config is empty, deleting")
		_ = os.Remove(configPath)
		return false
	}
	return true
}

func NewConfig() *Config {
	return &Config{
		Projects: []Project{},
	}
}

func LoadConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	configPath := filepath.Join(homeDir, ".wo", "wo_config.json")
	configStr, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(configStr, &config)
	if err != nil {
		return nil
	}
	return &config
}

func (c *Config) Save() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := filepath.Join(homeDir, ".wo", "wo_config.json")
	configStr, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(configPath, configStr, 0644)
}

func (c *Config) CreateProject(projectName string, projectDir string, commandStr string, args ...string) {
	if c.GetProject(projectName) != nil {
		fmt.Println("Project already exists")
		return
	}

	project := Project{
		Name: projectName,
		Dir:  projectDir,
		OpenFlow: OpenFlow{
			Timestamp: time.Now().Unix(),
			Commands:  []Command{CreateCommand(commandStr, args...)},
		},
		Timestamp: time.Now().Unix(),
	}
	c.AddProject(project)
	c.Save()
}

func (c *Config) Print() {
	for _, project := range c.Projects {
		fmt.Println(project.Name)
	}
}

func (c *Config) AddProject(project Project) {
	c.Projects = append(c.Projects, project)
}

func (c *Config) RemoveProject(project Project) {
	for i, p := range c.Projects {
		if p.Name == project.Name {
			c.Projects = append(c.Projects[:i], c.Projects[i+1:]...)
		}
	}
}

func (c *Config) GetProject(projectName string) *Project {
	for _, p := range c.Projects {
		if p.Name == projectName {
			return &p
		}
	}
	return nil
}

func (c *Config) GetProjectIndex(projectName string) int {
	for i, p := range c.Projects {
		if p.Name == projectName {
			return i
		}
	}
	return -1
}

func (c *Config) InsertCommand(projectName string, command Command) {
	projectIndex := c.GetProjectIndex(projectName)
	if projectIndex == -1 {
		fmt.Println("Project not found")
		return
	}
	project := c.Projects[projectIndex]
	project.OpenFlow.Commands = append(project.OpenFlow.Commands, command)
	c.Projects[projectIndex] = project
	c.Save()
	fmt.Println("Command added")
}

func (c *Config) RemoveCommand(projectName string, commandIndex int) {
	projectIndex := c.GetProjectIndex(projectName)
	if projectIndex == -1 {
		fmt.Println("Project not found")
		return
	}

	if commandIndex < 0 || commandIndex >= len(c.Projects[projectIndex].OpenFlow.Commands) {
		fmt.Println("Invalid command index, must be between 0 and", len(c.Projects[projectIndex].OpenFlow.Commands)-1)
		return
	}

	project := c.Projects[projectIndex]
	project.OpenFlow.Commands = append(project.OpenFlow.Commands[:commandIndex], project.OpenFlow.Commands[commandIndex+1:]...)
	c.Projects[projectIndex] = project
	c.Save()
	fmt.Println("Command removed")
}

func (c *Command) Run(dir string) {
	RunCommand(c.Command, dir, c.Args...)
}

func CreateCommand(commandStr string, args ...string) Command {
	return Command{
		Command: commandStr,
		Args:    args,
	}
}

func (o *OpenFlow) Run(dir string) {
	for _, command := range o.Commands {
		command.Run(dir)
	}
}

func (p *Project) Open() {
	p.OpenFlow.Run(p.Dir)
}
