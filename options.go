package gar

type Options struct {
	AnsibleRunnerBinary string // Name of the binary that will be used.
	AnsibleRunnerDir    string // The path to the directory where ansible-runner structure is defined.

	Playbook string // Name of playbook to run
	Ident    string // ansible-runner identifier directory name

	EnvVars map[string]string

	Logger Logger // Logger interface used
}

func (o Options) clone() Options {
	return Options{
		AnsibleRunnerBinary: o.AnsibleRunnerBinary,
		AnsibleRunnerDir:    o.AnsibleRunnerDir,
		Playbook:            o.Playbook,
		Ident:               o.Ident,
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

	return resultOptions, args
}
