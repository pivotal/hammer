package actions

import (
	"fmt"

	"github.com/pivotal/pcf/lockfile"
)

type OpenScripter struct{}

func NewOpenScripter() OpenScripter {
	return OpenScripter{}
}

func (b OpenScripter) Generate(data lockfile.Lockfile) []string {
	openCommandLines := []string{
		fmt.Sprintf(`open "%s"`, data.OpsManager.URL.String()),
		fmt.Sprintf(`echo "%s" | pbcopy`, data.OpsManager.Password),
	}

	return openCommandLines
}
