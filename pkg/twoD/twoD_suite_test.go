package twoD_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTwoD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TwoD Suite")
}
