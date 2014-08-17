package workbench_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	. "github.com/iancmcc/jig/workbench"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workbench", func() {

	var (
		dirname string
		bench   *Workbench
	)

	// Set up a temp directory
	BeforeEach(func() {
		dirname, _ = ioutil.TempDir("", "jig")
	})

	// Clean up the temp directory
	AfterEach(func() {
		if dirname != "" {
			syscall.Unlink(dirname)
		}
	})

	BeforeEach(func() {
		bench = NewWorkbench(dirname)
	})

	Context("creating a new workbench", func() {

		It("should respect the specified root directory", func() {
			Expect(bench.Root()).To(Equal(dirname))
			Expect(bench.SrcRoot()).To(Equal(filepath.Join(dirname, "src")))
			Expect(bench.BinDir()).To(Equal(filepath.Join(dirname, "bin")))
			Expect(bench.MetadataDir()).To(Equal(filepath.Join(dirname, ".jig")))
		})
	})

	verifyDirs := func() {
		st, err := os.Stat(bench.SrcRoot())
		Expect(st.IsDir()).To(BeTrue())
		Expect(err).To(BeNil())

		st, err = os.Stat(bench.BinDir())
		Expect(st.IsDir()).To(BeTrue())
		Expect(err).To(BeNil())

		st, err = os.Stat(bench.MetadataDir())
		Expect(st.IsDir()).To(BeTrue())
		Expect(err).To(BeNil())
	}

	Context("initializing a new workbench", func() {
		JustBeforeEach(func() {
			bench.Initialize()
		})

		It("should create the correct directories", func() {
			verifyDirs()
		})

		It("should create a valid plan", func() {
			Expect(bench.Plan()).ToNot(BeNil())
			st, err := os.Stat(filepath.Join(bench.MetadataDir(), "plan"))
			Expect(st).ToNot(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("existing directories", func() {
			BeforeEach(func() {
				os.MkdirAll(bench.MetadataDir(), 0755)
				os.MkdirAll(bench.SrcRoot(), 0755)
				os.MkdirAll(bench.BinDir(), 0755)
			})

			It("shouldn't fail if a directory exists", func() {
				verifyDirs()
			})
		})

	})

})
