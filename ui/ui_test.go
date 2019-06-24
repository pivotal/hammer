package ui_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/pivotal/pcf-cli/ui"
)

var _ = Describe("ui", func() {
	var ui UI

	BeforeEach(func() {
		ui = UI{
			Out: gbytes.NewBuffer(),
			Err: gbytes.NewBuffer(),
		}
	})

	Context("DisplayText", func() {
		It("prints text to the out buffer", func() {
			ui.DisplayText("Test text output")
			Expect(ui.Out).To(gbytes.Say("Test text output"))
		})
	})

	Context("DisplayError", func() {
		It("prints error text to the err buffer", func() {
			ui.DisplayError(fmt.Errorf("test error contents"))
			Expect(ui.Err).To(gbytes.Say("test error contents"))
		})
	})
})
