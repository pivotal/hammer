package scripting

import (
	"fmt"
	"os"
	"os/exec"
)

//go:generate counterfeiter . ScriptRunner

type ScriptRunner interface {
	RunScript(lines []string, prereqs []string, onlyWriteFile bool) error
}

type scriptRunner struct{}

func NewScriptRunner() ScriptRunner {
	return scriptRunner{}
}

func (s scriptRunner) RunScript(lines []string, prereqs []string, onlyWriteFile bool) error {
	path, err := WriteTempFile(lines...)
	if err != nil {
		return err
	}

	if onlyWriteFile {
		fmt.Println(path)
		return nil
	}

	if err := CheckPrereqs(prereqs); err != nil {
		return err
	}

	defer os.Remove(path)

	debug := "+x"
	if len(os.Getenv("DEBUG")) > 0 {
		debug = "-x"
	}

	command := exec.Command("/bin/bash", debug, path)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	return command.Run()
}
