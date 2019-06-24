package cf_test

import (
	"fmt"
	"net/url"

	"github.com/pivotal/pcf-cli/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf-cli/environment"
	"github.com/pivotal/pcf-cli/scripting/scriptingfakes"
)

var _ = Describe("cf login runner", func() {
	var (
		err           error
		cfLoginRunner cf.LoginRunner
		scriptRunner  *scriptingfakes.FakeScriptRunner

		data   environment.Config
		dryRun bool
	)

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

	It("runs the script with a cf login", func() {
		Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

		lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
		Expect(lines).To(ContainElement(`prods="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/staged/products)"`))
		Expect(lines).To(ContainElement(`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"`))
		Expect(lines).To(ContainElement(`cf login -a "api.sys.test-url.io" -u "$user" -p "$pass" --skip-ssl-validation`))

		Expect(prereqs).To(ConsistOf("jq", "om", "cf"))
		Expect(dryRun).To(Equal(true))
	})

	When("running the script succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("running the script errors", func() {
		BeforeEach(func() {
			scriptRunner.RunScriptReturns(fmt.Errorf("run-script-error"))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("run-script-error"))
		})
	})
})
