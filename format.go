package gar

import (
	"fmt"
	"os"
)

func listContains(l []string, s string) bool {
	for _, str := range l {
		if s == str {
			return true
		}
	}
	return false
}

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
