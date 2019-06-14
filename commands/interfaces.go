package commands

import "github.com/pivotal/pcf/environment"

//go:generate counterfeiter . EnvReader

type EnvReader interface {
	Read(targetConfigPath string) (environment.Config, error)
}

//go:generate counterfeiter . ToolRunner

type ToolRunner interface {
	Run(data environment.Config, dryRun bool, args ...string) error
}
