package scripting_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal/pcf-cli/scripting"
)

var _ = Describe("WriteTempFile", func() {
	It("returns an empty file if given no input lines", func() {
		path, err := WriteTempFile()
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(path)

		contents, err := ioutil.ReadFile(path)
		Expect(string(contents)).To(BeEmpty())
	})

	It("returns file containing input lines", func() {
		path, err := WriteTempFile("line1", "line2", "line3")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(path)

		contents, err := ioutil.ReadFile(path)
		Expect(string(contents)).To(Equal("line1\nline2\nline3\n"))
	})
})
