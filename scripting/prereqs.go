package scripting

import (
	"fmt"
	"os/exec"
)

func CheckPrereqs(prereqs []string) error {
	for _, v := range prereqs {
		_, err := exec.LookPath(v)
		if err != nil {
			return fmt.Errorf("Missing prerequisite '%s'. This must be installed first", v)
		}
	}

	return nil
}
