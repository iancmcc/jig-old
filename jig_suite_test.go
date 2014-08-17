package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jig Suite")
}
