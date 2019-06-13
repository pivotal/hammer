package testhelpers

import (
	"io/ioutil"
	"path"

	. "github.com/onsi/gomega"
)

func LoadFixture(name string) string {
	contents, err := ioutil.ReadFile(path.Join("fixtures", name))
	Expect(err).NotTo(HaveOccurred())
	return string(contents)
}
