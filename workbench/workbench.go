package workbench

import (
	"os"
	"path/filepath"

	"github.com/iancmcc/jig/plan"
)

// Workbench represents an entire Jig environment. It is rooted somewhere on
// the filesystem, and comprises a source root and bin dir.
type Workbench struct {
	root string     // Directory at the top of the bench
	plan *plan.Plan // The overall bench plan
}

// NewWorkbench creates a Workbench rooted at the specified directory.
func NewWorkbench(dirname string) *Workbench {
	return &Workbench{
		root: dirname,
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

// Initialize realizes the workbench on the filesystem.
func (b *Workbench) Initialize() error {
	if err := b.ensureDirectories(); err != nil {
		return err
	}
	if err := b.initializePlan(); err != nil {
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
	b.Plan().Repos[r.FullName] = r
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
	var p *plan.Plan
	filename := filepath.Join(b.MetadataDir(), "plan")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist yet; initialize with a new Plan
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		p = plan.NewPlan()
		p.ToJSON(f)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		p, err = plan.NewPlanFromJSON(f)
		if err != nil {
			return err
		}
	}
	if p.Repos == nil {
		p.Repos = make(map[string]*plan.Repo)
	}
	b.plan = p
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