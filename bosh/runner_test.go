package bosh_test

import (
	"fmt"
	"net"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/pcf-cli/bosh"
	"github.com/pivotal/pcf-cli/environment"
	"github.com/pivotal/pcf-cli/scripting/scriptingfakes"
)

var _ = Describe("bosh runner", func() {
	var (
		err          error
		boshRunner   bosh.Runner
		scriptRunner *scriptingfakes.FakeScriptRunner

		data     environment.Config
		dryRun   bool
		boshArgs []string
	)

	BeforeEach(func() {
		scriptRunner = new(scriptingfakes.FakeScriptRunner)

		url, _ := url.Parse("www.test-url.io")
		data = environment.Config{
			Name: "env-name",
			OpsManager: environment.OpsManager{
				PrivateKey: "private-key-contents",
				IP:         net.ParseIP("10.0.0.6"),
				URL:        *url,
				Username:   "username",
				Password:   "password",
			},
		}
		dryRun = true

		boshRunner = bosh.Runner{
			ScriptRunner: scriptRunner,
		}
	})

	JustBeforeEach(func() {
		err = boshRunner.Run(data, dryRun, boshArgs...)
	})

	When("no bosh args are passed to the bosh runner", func() {
		BeforeEach(func() {
			boshArgs = []string{}
		})

		It("invokes script runner with a series of bosh env var echos", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ContainElement(`echo "private-key-contents" >"$ssh_key_path"`))
			Expect(lines).To(ContainElement(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" ubuntu@"10.0.0.6" cat /var/tempest/workspaces/default/root_ca_certificate 1>${bosh_ca_path} 2>/dev/null`))
			Expect(lines).To(ContainElement(`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials)"`))
			Expect(lines).To(ContainElement(`bosh_proxy="BOSH_ALL_PROXY=ssh+socks5://ubuntu@10.0.0.6:22?private-key=${ssh_key_path}"`))
			Expect(lines).To(ContainElement(`echo "export BOSH_ENV_NAME=env-name"`))

			Expect(prereqs).To(ConsistOf("jq", "om", "ssh", "bosh"))
			Expect(dryRun).To(Equal(true))

			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("one or more bosh args are passed to the bosh runner", func() {
		BeforeEach(func() {
			boshArgs = []string{"arg1", "arg2", "arg3"}
			dryRun = false
		})

		It("invokes script runner with a bosh command", func() {
			Expect(scriptRunner.RunScriptCallCount()).To(Equal(1))

			lines, prereqs, dryRun := scriptRunner.RunScriptArgsForCall(0)
			Expect(lines).To(ContainElement(`echo "private-key-contents" >"$ssh_key_path"`))
			Expect(lines).To(ContainElement(`ssh -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -i "${ssh_key_path}" ubuntu@"10.0.0.6" cat /var/tempest/workspaces/default/root_ca_certificate 1>${bosh_ca_path} 2>/dev/null`))
			Expect(lines).To(ContainElement(`creds="$(om -t www.test-url.io -k -u username -p password curl -s -p /api/v0/deployed/director/credentials/bosh_commandline_credentials)"`))
			Expect(lines).To(ContainElement(`bosh_proxy="BOSH_ALL_PROXY=ssh+socks5://ubuntu@10.0.0.6:22?private-key=${ssh_key_path}"`))
			Expect(lines).To(ContainElement(`/usr/bin/env $bosh_client $bosh_env $bosh_secret $bosh_ca_cert $bosh_proxy bosh arg1 arg2 arg3`))

			Expect(prereqs).To(ConsistOf("jq", "om", "ssh", "bosh"))
			Expect(dryRun).To(Equal(false))

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
