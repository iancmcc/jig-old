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
	return nil
}

func (b *Workbench) initializePlan() error {
	filename := filepath.Join(b.MetadataDir(), "plan")
	f, err := os.Open(filename)
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