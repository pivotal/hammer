package scripting_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal/pcf-cli/scripting"
)

var _ = Describe("CheckPrereqs", func() {
	It("finds ls", func() {
		err := CheckPrereqs([]string{"ls"})
		Expect(err).NotTo(HaveOccurred())
	})

	It("finds multiple commands", func() {
		err := CheckPrereqs([]string{"ls", "cp"})
		Expect(err).NotTo(HaveOccurred())
	})

	It("complains when a command is not found", func() {
		err := CheckPrereqs([]string{"ls", "does-not-exist-fahfdslakfhsklfhsdgf", "cp"})
		Expect(err).To(MatchError("Missing prerequisite 'does-not-exist-fahfdslakfhsklfhsdgf'. This must be installed first"))
	})
})
