package commands

import "github.com/pivotal/pcf/lockfile"

//go:generate counterfeiter . ScriptRunner
type ScriptRunner interface {
	RunScript(lines []string, prereqs []string, onlyWriteFile bool) error
}

//go:generate counterfeiter . EnvReader
type EnvReader interface {
	Read(lockfilePath string) (lockfile.Lockfile, error)
}
