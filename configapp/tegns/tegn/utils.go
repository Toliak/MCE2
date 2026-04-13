package tegn

import (
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v6"
)

func defaultGitCloneOptions(additionalOpts ...func(*git.CloneOptions)) *git.CloneOptions {
	v := git.CloneOptions{
		Progress: os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		NoCheckout: true,
	}

	for _, opt := range additionalOpts {
		opt(&v)
	}

	return &v
}

func MkdirAllParent(path string) error {
	parentDir := filepath.Dir(path)
	return os.MkdirAll(parentDir, 0755)
}