package gar

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Command struct {
	Command    string            // The command to run
	Args       []string          // The args to pass to command
	WorkingDir string            // The directory where to run command
	Env        map[string]string // The environments variables added to command
	Logger     Logger            // The logger to use fo command's output
}

func generateCommand(options Options, args ...string) Command {
	return Command{
		Command:    options.AnsibleRunnerBinary,
		Args:       args,
		WorkingDir: options.AnsibleRunnerDir,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
}

func runAnsibleRunnerCommand(originalOptions Options, additionalArgs ...string) (string, error) {
	options, args := getCommonOptions(originalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.AnsibleRunnerBinary, args)
	return runCommandAndGetOutput(description, cmd)
}

func runCommandAndGetOutput(description string, cmd Command) (string, error) {
	cmd.Logger.Debug(description)
	output, err := runCommand(cmd)
	if err != nil {
		if output != nil {
			return output.Combined(), &ErrWithCmdOutput{err, output}
		}
		return "", &ErrWithCmdOutput{err, output}
	}

	return output.Combined(), nil
}

func runCommand(command Command) (*output, error) {
	command.Logger.Debug("Running command %s with args %s\n", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Env = formatEnvVars(command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	output, err := readStdoutAndStderr(command, stdout, stderr)
	if err != nil {
		return output, err
	}

	return output, cmd.Wait()
}

func readStdoutAndStderr(command Command, stdout, stderr io.ReadCloser) (*output, error) {
	out := newOutput()
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	var stdoutErr, stderrErr error
	go func() {
		defer wg.Done()
		stdoutErr = readData(command, false, stdoutReader, out.stdout)
	}()
	go func() {
		defer wg.Done()
		stderrErr = readData(command, true, stderrReader, out.stderr)
	}()
	wg.Wait()

	if stdoutErr != nil {
		return out, stdoutErr
	}
	if stderrErr != nil {
		return out, stderrErr
	}

	return out, nil
}

func readData(command Command, isStderr bool, reader *bufio.Reader, writer io.StringWriter) error {
	var line string
	var readErr error
	for {
		line, readErr = reader.ReadString('\n')

		// remove newline, our output is in a slice,
		// one element per line.
		line = strings.TrimSuffix(line, "\n")

		// only return early if the line does not have
		// any contents. We could have a line that does
		// not not have a newline before io.EOF, we still
		// need to add it to the output.
		if len(line) == 0 && readErr == io.EOF {
			break
		}

		if isStderr {
			command.Logger.Warn(line)
		} else {
			command.Logger.Info(line)
		}

		if _, err := writer.WriteString(line); err != nil {
			return err
		}

		if readErr != nil {
			break
		}
	}
	if readErr != io.EOF {
		return readErr
	}
	return nil
}
