package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	EnvReadError = errors.New("error reading environment vars")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, readDirErr := os.ReadDir(dir)
	if readDirErr != nil {
		return nil, fmt.Errorf(EnvReadError.Error()+": %w", readDirErr)
	}

	environment := make(Environment, len(files))
	for _, dirEntry := range files {
		file, openFileErr := os.Open(filepath.Join(dir, dirEntry.Name()))
		if openFileErr != nil {
			return nil, fmt.Errorf(EnvReadError.Error()+": %w", openFileErr)
		}
		defer file.Close()

		fileInfo, fileInfoErr := file.Stat()
		if fileInfoErr != nil {
			return nil, fmt.Errorf(EnvReadError.Error()+": %w", fileInfoErr)
		}

		if fileInfo.Size() == 0 {
			environment[dirEntry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		line := scanner.Text()

		environment[dirEntry.Name()] = EnvValue{Value: strings.TrimRight(line, "\x00\n")}
	}

	return environment, nil
}

func EmptyEnv() Environment {
	return make(map[string]EnvValue)
}
