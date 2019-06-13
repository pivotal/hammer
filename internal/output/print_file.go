package output

import (
	"fmt"
	"io/ioutil"
)

func WriteTempFile(lines ...string) (string, error) {
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	for _, l := range lines {
		if _, err = tempfile.Write([]byte(fmt.Sprintf("%s\n", l))); err != nil {
			return "", err
		}
	}

	if err = tempfile.Close(); err != nil {
		return "", err
	}

	return tempfile.Name(), nil
}
