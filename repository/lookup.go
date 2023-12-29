package repository

import (
	git "github.com/libgit2/git2go/v34"
)

type Commit struct {
	Commit    *git.Oid
	ShortId   string
	Tree      *git.Tree
	Message   string
	Author    *git.Signature
	Committer *git.Signature
	Parent    *git.Commit
}

type Tree struct {
	Tree    *git.Oid
	Entries []*git.TreeEntry
}

type Blob struct {
	Tree     *git.Oid
	ShortId  string
	IsBinary bool
	Contents []byte
}

type Tag struct {
	Tag    string
	Commit *Commit
}

// LookupCommit - Lookup for commit oid
func LookupCommit(repositoryName, commitId string) (*Commit, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	oid, err := git.NewOid(commitId)
	if err != nil {
		return nil, handleGitError(err, "unable to parse oid")
	}

	commit, err := repository.LookupCommit(oid)
	if err != nil {
		return nil, handleGitError(err, "unable to lookup commit")
	}

	return GetCommit(commit)
}

// LookupTree - Lookup for tree oid
func LookupTree(repositoryName, treeId string) (*Tree, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	oid, err := git.NewOid(treeId)
	if err != nil {
		return nil, handleGitError(err, "unable to parse oid")
	}

	tree, err := repository.LookupTree(oid)
	if err != nil {
		return nil, handleGitError(err, "unable to lookup commit")
	}

	return GetTree(tree)
}

// LookupBlob - Lookup for blob oid
func LookupBlob(repositoryName, blobId string) (*Blob, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	oid, err := git.NewOid(blobId)
	if err != nil {
		return nil, handleGitError(err, "unable to parse oid")
	}

	blob, err := repository.LookupBlob(oid)
	if err != nil {
		return nil, handleGitError(err, "unable to lookup commit")
	}

	return GetBlob(blob)
}

// LookupTag - Lookup for tag name
func LookupTag(repositoryName, tagName string) (*Tag, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	object, err := repository.RevparseSingle(tagName)
	if err != nil {
		return nil, handleGitError(err, "unable to rev parse tag")
	}

	tag, err := repository.LookupTag(object.Id())
	if err != nil {
		return nil, handleGitError(err, "unable to lookup tag")
	}

	commit, err := repository.LookupCommit(tag.TargetId())
	if err != nil {
		return nil, handleGitError(err, "unable to lookup commit")
	}

	tagCommit, err := GetCommit(commit)
	if err != nil {
		return nil, handleGitError(err, "unable to get tag commit")
	}

	return &Tag{
		Tag:    tagName,
		Commit: tagCommit,
	}, nil
}

// GetCommit = Get commit information
func GetCommit(commit *git.Commit) (*Commit, error) {
	tree, err := commit.Tree()
	if err != nil {
		return nil, handleGitError(err, "unable to get commit tree")
	}

	shortId, err := commit.ShortId()
	if err != nil {
		return nil, handleGitError(err, "unable to get short id")
	}

	return &Commit{
		Commit:    commit.Id(),
		ShortId:   shortId,
		Tree:      tree,
		Message:   commit.Message(),
		Author:    commit.Author(),
		Committer: commit.Committer(),
		Parent:    commit.Parent(0),
	}, nil
}

// GetTree - Get tree entries
func GetTree(tree *git.Tree) (*Tree, error) {
	var entries []*git.TreeEntry
	err := tree.Walk(func(s string, entry *git.TreeEntry) error {
		entries = append(entries, entry)
		return nil
	})

	if err != nil {
		return nil, handleGitError(err, "unable to walk throught the tree")
	}

	return &Tree{
		Tree:    tree.Id(),
		Entries: entries,
	}, nil
}

// GetBlob - Get blob contents
func GetBlob(blob *git.Blob) (*Blob, error) {
	shortId, err := blob.ShortId()
	if err != nil {
		return nil, handleGitError(err, "unable to get short id")
	}
	return &Blob{
		Tree:     blob.Id(),
		ShortId:  shortId,
		IsBinary: blob.IsBinary(),
		Contents: blob.Contents(),
	}, nil
}
