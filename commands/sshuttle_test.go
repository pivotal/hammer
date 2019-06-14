package commands_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes "github.com/pivotal/pcf/commands/commandsfakes"
	"github.com/pivotal/pcf/lockfile"

	. "github.com/pivotal/pcf/commands"
)

var _ = Describe("sshuttle command", func() {
	var err error
	var command *SshuttleCommand

	var envReader *fakes.FakeEnvReader
	var sshuttleRunner *fakes.FakeToolRunner
	var args []string
	var dryRun bool

	BeforeEach(func() {
		envReader = new(fakes.FakeEnvReader)
		sshuttleRunner = new(fakes.FakeToolRunner)
		args = []string{"arg1", "arg2"}
		dryRun = true
	})

	JustBeforeEach(func() {
		command = &SshuttleCommand{
			Env:            envReader,
			SshuttleRunner: sshuttleRunner,
			File:           dryRun,
		}

		err = command.Execute(args)
	})

	When("envReader returns an error", func() {
		BeforeEach(func() {
			envReader.ReadReturns(lockfile.Lockfile{}, fmt.Errorf("env-reader-error"))
		})

		It("doesn't call sshuttleRunner", func() {
			Expect(sshuttleRunner.RunCallCount()).To(Equal(0))
		})

		It("propagates the error", func() {
			Expect(err).To(MatchError("env-reader-error"))
		})
	})

	When("envReader succeeds", func() {
		BeforeEach(func() {
			envReader.ReadReturns(lockfile.Lockfile{Name: "env-name"}, nil)
		})

		It("passes env data and dry run flag but not args to sshuttleRunner", func() {
			Expect(sshuttleRunner.RunCallCount()).To(Equal(1))

			environmentConfig, dryRun, args := sshuttleRunner.RunArgsForCall(0)
			Expect(environmentConfig).To(BeEquivalentTo(lockfile.Lockfile{Name: "env-name"}))
			Expect(dryRun).To(BeTrue())
			Expect(args).To(HaveLen(0))
		})

		When("sshuttleRunner succeeds", func() {
			BeforeEach(func() {
				sshuttleRunner.RunReturns(nil)
			})

			It("doesn't error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("sshuttleRunner returns an error", func() {
			BeforeEach(func() {
				sshuttleRunner.RunReturns(fmt.Errorf("sshuttle-runnner-error"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("sshuttle-runnner-error"))
			})
		})
	})
})
