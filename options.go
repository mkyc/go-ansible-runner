package gar

type LogsLevel int

const (
	L1 LogsLevel = iota
	L2
	L3
	L4
	L5
	L6
)

type Options struct {
	AnsibleRunnerBinary string // Name of the binary that will be used.
	AnsibleRunnerDir    string // The path to the directory where ansible-runner structure is defined.

	Playbook string // Name of playbook to run
	Ident    string // ansible-runner identifier directory name

	EnvVars   map[string]string
	LogsLevel LogsLevel

	Logger Logger // Logger interface used
}

func (o Options) clone() Options {
	return Options{
		AnsibleRunnerBinary: o.AnsibleRunnerBinary,
		AnsibleRunnerDir:    o.AnsibleRunnerDir,
		Playbook:            o.Playbook,
		Ident:               o.Ident,
		EnvVars:             o.EnvVars,
		LogsLevel:           o.LogsLevel,
		Logger:              o.Logger,
	}
}

func getCommonOptions(options Options, args ...string) (Options, []string) {
	resultOptions := options.clone()
	if options.AnsibleRunnerBinary == "" {
		resultOptions.AnsibleRunnerBinary = "ansible-runner"
	}
	if options.AnsibleRunnerDir == "" {
		resultOptions.AnsibleRunnerDir = "."
	}
	if options.Ident != "" && !listContains(args, "--ident") {
		args = append(args, "--ident", options.Ident)
	}
	if options.Playbook != "" && !listContains(args, "--playbook") {
		args = append(args, "--playbook", options.Playbook)
	}
	if options.Logger == nil {
		resultOptions.Logger = DefaultLogger{}
	}
	switch options.LogsLevel {
	case L1:
		args = append(args, "-v")
	case L2:
		args = append(args, "-vv")
	case L3:
		args = append(args, "-vvv")
	case L4:
		args = append(args, "-vvvv")
	case L5:
		args = append(args, "-vvvvv")
	case L6:
		args = append(args, "-vvvvv", "--debug")
	}
	return resultOptions, args
}
