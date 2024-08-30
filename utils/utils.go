package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

var (
	logger = log.Default()
)

type ExecConfig struct {
	Command     string
	Args        []string
	WorkingDir  string
	Environment map[string]string
}

func CheckErr(err error) {
	if err != nil {
		logger.Fatalf(fmt.Sprintf("ERROR: %s", err))
	}
}

func Exec(config ExecConfig) (string, error) {
	cmd := exec.Command(config.Command, config.Args...)

	var output = &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output

	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}

	if len(config.Environment) > 0 {
		for key, val := range config.Environment {
			cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", key, val))
		}
	}

	err := cmd.Run()
	if err != nil {
		return output.String(), fmt.Errorf("there was an error executing the command: %s %s. ERROR: %s", config.Command, config.Args, err)
	}

	return output.String(), nil
}
