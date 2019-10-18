package control_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestControl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Control Suite")
}
