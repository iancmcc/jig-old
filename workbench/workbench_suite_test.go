package workbench_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkbench(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workbench Suite")
}
