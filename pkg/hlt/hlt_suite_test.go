package hlt_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHlt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hlt Suite")
}
