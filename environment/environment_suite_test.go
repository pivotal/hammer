package environment_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Environment Suite")
}
