package internal

type StageType string

const (
	ENVIRONMENTS StageType = "ENVIRONMENTS"
	COMMANDS     StageType = "COMMANDS"
	ENTRIES      StageType = "ENTRIES"
)

type Project struct {
	Name   string
	Stages []Stage
}

type Stage struct {
	Type        StageType
	Environment map[string]string `json:"Environment,omitempty"`
	Commands    []string          `json:"Commands,omitempty"`
	Entries     []Entry           `json:"Entries,omitempty"`
}

type Entry struct {
	IsDir   bool
	Encoded bool
	Path    string
	Content string `json:"Content,omitempty"`
}
