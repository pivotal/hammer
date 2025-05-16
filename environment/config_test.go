/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package environment_test

import (
	"net"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/pivotal/hammer/environment"
)

var _ = Describe("Config", func() {
	AfterEach(func() {
		os.Unsetenv("HAMMER_TARGET_CONFIG")
	})

	Describe("FromFile", func() {
		It("reads data from a json file", func() {
			env, err := FromFile(path.Join("fixtures", "lemon.json"), "")

			Expect(err).NotTo(HaveOccurred())
			checkMatchLemon(env)
		})

		It("reads data from a yaml file", func() {
			env, err := FromFile(path.Join("fixtures", "lemon.yaml"), "")

			Expect(err).NotTo(HaveOccurred())
			checkMatchLemon(env)
		})

		It("reads data from a config file that does not contain subnets, CIDRs, AZs or version", func() {
			env, err := FromFile(path.Join("fixtures", "reduced.json"), "")

			Expect(err).NotTo(HaveOccurred())
			checkMatchReduced(env)
		})

		When("a non-matching environment name is specified", func() {
			It("errors with a helpful message", func() {
				env, err := FromFile(path.Join("fixtures", "lemon.json"), "non-matching-environment")

				Expect(err).To(MatchError("Environment name 'non-matching-environment' specified but does not match environment in config"))
				Expect(env).To(BeEquivalentTo(Config{}))
			})
		})

		When("a file with multiple environment configs is specified", func() {
			It("selects the first config if an environment name is not specified", func() {
				env, err := FromFile(path.Join("fixtures", "multiple.json"), "")

				Expect(err).NotTo(HaveOccurred())
				checkMatchLemon(env)
			})

			It("selects the appropriate config if a matching environment name is specified", func() {
				env, err := FromFile(path.Join("fixtures", "multiple.json"), "reduced-config")

				Expect(err).NotTo(HaveOccurred())
				checkMatchReduced(env)
			})

			It("errors if a non-matching environment name is specified", func() {
				env, err := FromFile(path.Join("fixtures", "multiple.json"), "non-existent-config")

				Expect(err).To(MatchError("Environment name 'non-existent-config' specified but does not match environment in config"))
				Expect(env).To(BeEquivalentTo(Config{}))
			})
		})

		When("a file with an empty array is specified", func() {
			It("errors", func() {
				env, err := FromFile(path.Join("fixtures", "empty.json"), "")

				Expect(err).To(MatchError("Target config is an empty array"))
				Expect(env).To(BeEquivalentTo(Config{}))
			})
		})
	})
})

func mustParseURL(u string) url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return *url
}

func checkMatchLemon(e Config) {
	Expect(e).To(MatchAllFields(Fields{
		"Name":          Equal("lemon"),
		"Version":       Equal(*version.Must(version.NewVersion("1.11"))),
		"CFDomain":      Equal("sys.lemon.cf-app.com"),
		"AppsDomain":    Equal("apps.lemon.cf-app.com"),
		"PasSubnet":     Equal("lemon-pas-subnet"),
		"ServiceSubnet": Equal("lemon-services-subnet"),
		"AZs":           Equal([]string{"us-central1-f", "us-central1-a", "us-central1-c"}),
		"PKSApi": MatchFields(IgnoreExtras, Fields{
			"Username": Equal("pivotalcf"),
			"Password": Equal("fakePassword"),
			"URL":      Equal(mustParseURL("https://api.pks.lemon-lemon.cf-app.com")),
		}),
		"OpsManager": MatchAllFields(Fields{
			"Username":     Equal("pivotalcf"),
			"Password":     Equal("fakePassword"),
			"ClientID":     Equal("fakeClientID"),
			"ClientSecret": Equal("fakeClientSecret"),
			"URL":          Equal(mustParseURL("https://pcf.lemon.cf-app.com")),
			"IP":           Equal(net.ParseIP("35.225.148.133")),
			"SshUser":      Equal("ubuntu"),
			"PrivateKey":   ContainSubstring("BEGIN RSA"),
		}),
	}))
}

func checkMatchReduced(e Config) {
	Expect(e).To(MatchFields(IgnoreExtras, Fields{
		"Name":     Equal("reduced-config"),
		"CFDomain": Equal("sys.reduced-config.cf-app.com"),
		"PKSApi": MatchFields(IgnoreExtras, Fields{
			"URL": Equal(mustParseURL("https://api.pks.reduced-config.cf-app.com")),
		}),
		"OpsManager": MatchFields(IgnoreExtras, Fields{
			"Username":   Equal("pivotalcf"),
			"Password":   Equal("fakePassword"),
			"URL":        Equal(mustParseURL("https://pcf.reduced-config.cf-app.com")),
			"IP":         Equal(net.ParseIP("35.225.148.133")),
			"SshUser":    Equal("ubuntu"),
			"PrivateKey": ContainSubstring("BEGIN RSA"),
		}),
	}))
}
