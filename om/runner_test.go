package om_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/pcf-cli/om"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf-cli/environment"
	"github.com/pivotal/pcf-cli/scripting/scriptingfakes"
)

var _ = Describe("om runner", func() {

	var err error
	var omRunner om.Runner
	var scriptRunner *scriptingfakes.FakeScriptRunner

	var data environment.Config
	var dryRun bool
	var omArgs []string

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			OpsManager: environment.OpsManager{
				URL:      *url,
				Username: "username",
				Password: "password",
			},
		}
		dryRun = true

		omRunner = om.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = omRunner.Run(data, dryRun, omArgs...)
	})

	Context("run", func() {
		When("no om args are passed to the runner", func() {
			BeforeEach(func() {
				omArgs = []string{}
			})

			It("invokes script runner with a series of om env var echos", func() {
				Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

				lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
				Expect(lines).To(ConsistOf(
					`echo "export OM_TARGET=www.test-url.io"`,
					`echo "export OM_USERNAME=username"`,
					`echo "export OM_PASSWORD=password"`,
				))

				Expect(prereqs).To(Equal([]string{"om"}))
				Expect(dryRun).To(Equal(true))
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("one or more om args are passed to the runner", func() {
			BeforeEach(func() {
				omArgs = []string{"arg1", "arg2", "arg3"}
				dryRun = false
			})

			It("invokes script runner with an om command", func() {
				Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

				lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
				Expect(lines).To(ConsistOf(`om -t 'www.test-url.io' -k -u 'username' -p 'password' 'arg1' 'arg2' 'arg3'`))

				Expect(prereqs).To(ConsistOf("om"))
				Expect(dryRun).To(Equal(false))
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("script runner run script errors", func() {
			BeforeEach(func() {
				scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("run-script-error"))
			})
		})
	})
})
