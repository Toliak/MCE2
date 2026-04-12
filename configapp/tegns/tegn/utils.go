package tegn

import (
	"os"

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
