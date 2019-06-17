package cf_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/pcf/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf/environment"
	"github.com/pivotal/pcf/scripting/scriptingfakes"
)

var _ = Describe("cf login runner", func() {

	var err error
	var cfLoginRunner cf.LoginRunner
	var scriptRunner *scriptingfakes.FakeScriptRunner

	var data environment.Config
	var dryRun bool

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			CFDomain: "sys.test-url.io",
			OpsManager: environment.OpsManager{
				URL:      *url,
				Username: "username",
				Password: "password",
			},
		}
		dryRun = true

		cfLoginRunner = cf.LoginRunner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = cfLoginRunner.Run(data, dryRun)
	})

	Context("run", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(nil)
		})

		It("invokes script runner with a cf login", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ContainElement(`prods="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/staged/products)"`))
			Expect(lines).To(ContainElement(`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"`))
			Expect(lines).To(ContainElement(`cf login -a "api.sys.test-url.io" -u "$user" -p "$pass" --skip-ssl-validation`))

			Expect(prereqs).To(ConsistOf("jq", "om", "cf"))
			Expect(dryRun).To(Equal(true))
		})

		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
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
