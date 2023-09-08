package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdInfo []string, env Environment) (returnCode int) {
	setEnvVars(env)

	command := getCommand(cmdInfo)
	command.Run()

	return command.ProcessState.ExitCode()
}

func setEnvVars(env Environment) {
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
			continue
		}

		os.Setenv(name, value.Value)
	}
}

func getCommand(cmd []string) *exec.Cmd {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command
}
