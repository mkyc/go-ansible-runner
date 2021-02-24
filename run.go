package gar

func Run(options Options) (string, error) {
	return runAnsibleRunnerCommand(options, "run", ".")
}
