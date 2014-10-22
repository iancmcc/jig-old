package workbench

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/iancmcc/jig/plan"
	"github.com/iancmcc/jig/vcs"
)

// Workbench represents an entire Jig environment. It is rooted somewhere on
// the filesystem, and comprises a source root and bin dir.
type Workbench struct {
	root string     // Directory at the top of the bench
	plan *plan.Plan // The overall bench plan
}

// NewWorkbench creates a Workbench rooted at the specified directory.
func NewWorkbench(dirname string, plan *plan.Plan) *Workbench {
	return &Workbench{
		root: dirname,
		plan: plan,
	}
}

// Root returns the directory at which the Workbench is rooted.
func (b *Workbench) Root() string {
	return b.root
}

// SrcRoot returns the Workbench's source directory.
func (b *Workbench) SrcRoot() string {
	return filepath.Join(b.root, "src")
}

// BinDir returns the Workbench's bin directory.
func (b *Workbench) BinDir() string {
	return filepath.Join(b.root, "bin")
}

// MetadataDir returns the Workbench's metadata directory.
func (b *Workbench) MetadataDir() string {
	return filepath.Join(b.root, ".jig")
}

// Realize realizes the workbench on the filesystem.
func (b *Workbench) Realize() error {
	if err := b.ensureDirectories(); err != nil {
		return err
	}
	if err := b.initializePlan(); err != nil {
		return err
	}
	if err := b.initializeRepos(); err != nil {
		return err
	}
	return nil
}

// Plan returns the bench's current plan
func (b *Workbench) Plan() *plan.Plan {
	return b.plan
}

// AddRepository adds a repository to the workbench and saves the plan.
func (b *Workbench) AddRepository(r *plan.Repo) error {
	b.Plan().Repos[r.URI] = r
	return b.save()
}

func (b *Workbench) save() error {
	filename := filepath.Join(b.MetadataDir(), "plan")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b.Plan().ToJSON(f)
	return nil
}

func (b *Workbench) initializePlan() error {
	filename := filepath.Join(b.MetadataDir(), "plan")
	if b.plan == nil {
		b.plan = plan.NewPlan()
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist yet; initialize
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		b.plan.ToJSON(f)
	}
	return nil
}

func (b *Workbench) initializeRepos() error {
	var wg sync.WaitGroup
	bank := vcs.NewProgressBarBank()
	for name, r := range b.Plan().Repos {
		if name != "" {
			wg.Add(1)
			go func(name string, r *plan.Repo) {
				defer wg.Done()
				srcrepo, _ := vcs.NewSourceRepository(r, b.SrcRoot(), bank)
				if err := srcrepo.Create(); err != nil {
					fmt.Printf("%+v\n", err)
				}
			}(name, r)
		}
	}
	wg.Wait()
	return nil
}

func (b *Workbench) ensureDirectories() error {
	if err := os.MkdirAll(b.MetadataDir(), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(b.SrcRoot(), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(b.BinDir(), 0755); err != nil {
		return err
	}
	return nil
}

func FindNearestBench(dir string) *Workbench {
	for {
		metadir := path.Join(dir, ".jig")
		fmt.Println("Looking for", metadir)
		if _, err := os.Stat(metadir); os.IsNotExist(err) {
			dir = filepath.Dir(dir)
			fmt.Println("New dir", dir)
			if dir == "." || dir == "/" {
				break
			}
			continue
		}
		planfile := path.Join(metadir, "plan")
		f, err := os.Open(planfile)
		if err != nil {
			break
		}
		defer f.Close()
		p, err := plan.NewPlanFromJSON(f)
		return NewWorkbench(dir, p)
	}
	return nil
}