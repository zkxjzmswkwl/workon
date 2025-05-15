package wo

type Command struct {
	Command string
	Args    []string
}

type OpenFlow struct {
	Timestamp int64
	Commands  []Command
}

type Project struct {
	Timestamp int64
	Name      string
	Dir       string
	OpenFlow  OpenFlow
}

type Config struct {
	Projects []Project
}
