package plan_test

import (
	"fmt"
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
				"git@github.com:iancmcc/jig": {
					"ref": "develop"
				},
				"git@github.com:iancmcc/dotfiles": {
					"type": "subversion"
				}
			}
		}`))
	})

	Describe("loading from JSON", func() {

		Context("when the JSON parses successfully", func() {
			It("should populate fields properly", func() {
				repo := plan.Repos["git@github.com:iancmcc/jig"]
				fmt.Printf("%+v", plan)
				Expect(repo.RefSpec).To(Equal("develop"))
				Expect(repo.FullName).To(Equal("github.com/iancmcc/jig"))
			})
			It("should autodetect VCS type if possible", func() {
				repo := plan.Repos["git@github.com:iancmcc/jig"]
				Expect(repo.VCSType()).To(Equal(GIT))
			})
			It("should respect user-specified VCS type", func() {
				repo := plan.Repos["git@github.com:iancmcc/dotfiles"]
				Expect(repo.VCSType()).To(Equal(SVN))
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
