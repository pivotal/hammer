package environment_test

import (
	"net"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/pivotal/pcf-cli/environment"
)

var _ = Describe("Config", func() {

	AfterEach(func() {
		os.Unsetenv("TARGET_ENVIRONMENT_CONFIG")
	})

	Describe("FromFile", func() {
		It("reads data from a file", func() {
			env, err := FromFile(path.Join("fixtures", "lemon.json"))

			Expect(err).NotTo(HaveOccurred())
			checkMatchLemon(env)
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

func mustParseCIDR(c string) net.IPNet {
	_, cidr, err := net.ParseCIDR(c)
	if err != nil {
		panic(err)
	}
	return *cidr
}

func checkMatchLemon(e Config) {
	Expect(e).To(MatchAllFields(Fields{
		"Name":          Equal("lemon"),
		"Version":       Equal(*version.Must(version.NewVersion("1.11"))),
		"CFDomain":      Equal("sys.lemon.cf-app.com"),
		"AppsDomain":    Equal("apps.lemon.cf-app.com"),
		"PasCIDR":       Equal(mustParseCIDR("10.0.4.0/24")),
		"PasSubnet":     Equal("lemon-pas-subnet"),
		"ServicesCIDR":  Equal(mustParseCIDR("10.0.8.0/24")),
		"ServiceSubnet": Equal("lemon-services-subnet"),
		"AZs":           Equal([]string{"us-central1-f", "us-central1-a", "us-central1-c"}),
		"OpsManager": MatchAllFields(Fields{
			"Username":   Equal("pivotalcf"),
			"Password":   Equal("fakePassword"),
			"URL":        Equal(mustParseURL("https://pcf.lemon.cf-app.com")),
			"IP":         Equal(net.ParseIP("35.225.148.133")),
			"CIDR":       Equal(mustParseCIDR("10.0.0.0/24")),
			"PrivateKey": ContainSubstring("BEGIN RSA"),
		}),
	}))
}
