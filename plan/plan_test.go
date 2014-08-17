package plan_test

import (
	"strings"

	. "github.com/iancmcc/jig/plan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {

	var (
		plan *Plan
		err  error
	)

	BeforeEach(func() {
		plan, err = NewPlanFromJSON(strings.NewReader(`{
			"repos": {
				"github.com/iancmcc/jig": {
				}
			}
		}`))
	})

	Describe("loading from JSON", func() {

		Context("when the JSON parses successfully", func() {
			It("should populate fields properly", func() {
				Expect(plan.Repos).To(HaveLen(1))
				for name, repo := range plan.Repos {
					Expect(name).To(Equal("github.com/iancmcc/jig"))
					Expect(repo.FullName).To(Equal(name))
					break
				}
			})
		})
		Context("when the JSON fails to parse", func() {
			BeforeEach(func() {
				plan, err = NewPlanFromJSON(strings.NewReader(`{
					"ian": 123invalid
				}`))
			})
			It("should return the zero value for the plan", func() {
				Expect(plan).To(BeZero())
			})
			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
