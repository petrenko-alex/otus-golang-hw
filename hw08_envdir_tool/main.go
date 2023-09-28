package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]
	cmdAndArgs := os.Args[2:]

	env := EmptyEnv()
	if envDir != "" {
		var readDirErr error
		env, readDirErr = ReadDir(envDir)
		if readDirErr != nil {
			fmt.Println(readDirErr)
			os.Exit(1)
		}
	}

	exitCode, err := RunCmd(cmdAndArgs, env)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitCode)
}
