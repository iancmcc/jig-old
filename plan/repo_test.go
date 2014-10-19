package plan_test

import (
	. "github.com/iancmcc/jig/plan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {

	var (
		repo Repo
	)

	BeforeEach(func() {
		repo, _ = NewRepo("git@github.com:iancmcc/jig")
	})

	Describe("parsing the name", func() {
		Context("the name is well formed", func() {

			It("should parse the elements", func() {
				Expect(repo.Registry()).To(Equal("github.com"))
				Expect(repo.Owner()).To(Equal("iancmcc"))
				Expect(repo.Repository()).To(Equal("jig"))
			})

			It("should detect that github.com is git repos", func() {
				Expect(repo.VCSType()).To(Equal(GIT))
			})

			Context("the VCS type is specified", func() {
				BeforeEach(func() {
					repo.UserType = "subversion"
				})

				It("should use the specified VCS type", func() {
					Expect(repo.VCSType()).To(Equal(SVN))
				})
			})

		})

		Context("the name is not well formed", func() {
			BeforeEach(func() {
				repo.FullName = "github.com"
			})
			It("should error when parsing registry", func() {
				registry, err := repo.Registry()
				Expect(registry).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
			It("should error when parsing owner", func() {
				owner, err := repo.Owner()
				Expect(owner).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
			It("should error when parsing repository", func() {
				repository, err := repo.Repository()
				Expect(repository).To(BeZero())
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
