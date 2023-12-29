package repository

import (
	"errors"
	"fmt"

	git "github.com/libgit2/git2go/v34"
)

var (
	NotFoundError = errors.New("not found")
)

func handleGitError(err error, message string) error {
	var gitError *git.GitError
	isGitError := errors.As(err, &gitError)
	if isGitError {
		if gitError.Code == git.ErrorCodeNotFound {
			return NotFoundError
		}
	}
	return fmt.Errorf("%s %w", message, err)
}
