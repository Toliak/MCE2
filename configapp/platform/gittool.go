package platform

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type GitClone struct {
	// Git repository URL
	URL string

	// Filesystem path
	Path string

	// Remote branch.
	// Leave empty to use the remote's default branch
	Branch string

	// Performs a shallow clone with depth 1 if true
	ShallowMode bool
}

type GitCloneOption func(*GitClone)

// WithBranch sets the branch to clone.
func WithGitBranch(branch string) GitCloneOption {
	return func(o *GitClone) {
		o.Branch = branch
	}
}

// WithShallowMode enables or disables shallow cloning.
func WithGitShallowMode(shallow bool) GitCloneOption {
	return func(o *GitClone) {
		o.ShallowMode = shallow
	}
}

// NewGitClone creates a new GitCloneOptions instance using functional options.
func NewGitClone(url, path string, opts ...GitCloneOption) *GitClone {
	// Set defaults
	g := &GitClone{
		URL:        url,
		Path:       path,
		Branch:     "", // Default: use remote's default
		ShallowMode: false, // Default: full clone
	}

	// Apply provided options
	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Clone performs the git clone operation based on the configured options.
func (g *GitClone) Clone() error {
	if err := MkdirAllParent(g.Path); err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", g.Path, err)
	}

	cloneOptions := g.buildCloneOptions()

	fmt.Printf("Repository clone '%s' into '%s'\n", g.URL, g.Path)
	repo, err := git.PlainClone(g.Path, cloneOptions)
	if err != nil {
		return fmt.Errorf("PlainClone error: %w", err)
	}

	// If a specific branch was requested, checkout that branch.
	// PlainClone with specific Branch already fetches that branch's tip, so checkout should work.
	if g.Branch != "" {
		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("Worktree error: %w", err)
		}

		branchRefName := plumbing.NewBranchReferenceName(g.Branch)
		branchCoOpts := git.CheckoutOptions{
			Branch: branchRefName,
		}
		if err = w.Checkout(&branchCoOpts); err != nil {
			return fmt.Errorf("Checkout error for branch '%s': %w", g.Branch, err)
		}
	}

	return nil
}

// buildCloneOptions constructs the *git.CloneOptions based on the GitCloneOptions settings.
func (g *GitClone) buildCloneOptions() *git.CloneOptions {
	opts := &git.CloneOptions{
		URL:               g.URL,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if g.Branch != "" {
		opts.ReferenceName = plumbing.NewBranchReferenceName(g.Branch)
	}
	if g.ShallowMode {
		opts.SingleBranch = true
		opts.Depth = 1
	}

	return opts
}