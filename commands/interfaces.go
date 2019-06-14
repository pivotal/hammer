package commands

import "github.com/pivotal/pcf/lockfile"

//go:generate counterfeiter . EnvReader

type EnvReader interface {
	Read(lockfilePath string) (lockfile.Lockfile, error)
}

//go:generate counterfeiter . ToolRunner

type ToolRunner interface {
	Run(data lockfile.Lockfile, dryRun bool, args ...string) error
}
