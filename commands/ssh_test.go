package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf/commands/commandsfakes"
	"github.com/pivotal/pcf/lockfile"

	. "github.com/pivotal/pcf/commands"
)

var _ = Describe("ssh command", func() {
	var err error
	var command *SSHCommand

	var envReader *fakes.FakeEnvReader
	var sshRunner *fakes.FakeToolRunner
	var args []string
	var dryRun bool

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		sshRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}
		dryRun = true
	})

	JustBeforeEach(func() {
		command = &SSHCommand{
			Env:       envReader,
			SSHRunner: sshRunner,
			File:      dryRun,
		}

		err = command.Execute(args)
	})

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(lockfile.Lockfile{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call sshRunner", func() {
			Expect(sshRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(lockfile.Lockfile{Name: "env-name"}, nil)
		})

		It("passes env data and dry run flag but not args to sshRunner", func() {
			Expect(sshRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := sshRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(lockfile.Lockfile{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
		})

		When("sshRunner succeeds", func() {
			BeforeEach(func() {
				sshRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("sshRunner returns an error", func() {
			BeforeEach(func() {
				sshRunner.RunReturns(fmt.Errorf("ssh-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("ssh-runnner-error"))
			})
		})
	})
})
