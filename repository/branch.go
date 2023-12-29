package repository

import (
	git "github.com/libgit2/git2go/v34"
)

type Branch struct {
	Branch string
	Commit *git.Oid
	Tree   *git.Oid
	Parent *git.Oid
}

// ListRepositoryBranches - list all branches in the repository
func ListRepositoryBranches(repositoryName string) ([]Branch, error) {
	var branches []Branch
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	iterator, err := repository.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return nil, handleGitError(err, "unable to create branch iterator")
	}

	err = iterator.ForEach(func(b *git.Branch, bt git.BranchType) error {
		branchName, err := b.Name()
		if err != nil {
			return err
		}

		branch, err := GetBranch(repositoryName, branchName)
		if err != nil {
			return err
		}

		branches = append(branches, *branch)
		return nil
	})
	if err != nil {
		return nil, handleGitError(err, "unable to use branch iterator")
	}

	return branches, nil
}

// GetBranch - Get branch commit, tree and parent
func GetBranch(repositoryName, branchName string) (*Branch, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	branch, err := repository.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, handleGitError(err, "unable to lookup branch")
	}

	commit, err := repository.LookupCommit(branch.Target())
	if err != nil {
		return nil, handleGitError(err, "unable to lookup commit")
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, handleGitError(err, "unable to get commit tree")
	}

	parent := commit.ParentId(0)
	return &Branch{Branch: branchName, Commit: commit.Id(), Tree: tree.Id(), Parent: parent}, nil
}
