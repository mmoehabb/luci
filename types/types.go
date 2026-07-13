package types

type Config struct {
	Title       string
	Description string
	Bash        ShellConfig
	Zshell      ShellConfig
	Powershell  ShellConfig
	Wildcard    ShellConfig `toml:"*"`
}

type ShellConfig map[string]any // any: Action or ActionRecord

type ActionRecord map[string]Action

type Action any // VagueAction or AnnotatedAction

type VagueAction interface {
	string | []string
}

type AnnotatedAction struct {
	Title       string
	Description string
	Value       any
}

type ShellType int

const (
	Bash ShellType = iota
	Zshell
	Powershell
	Unknown
)
