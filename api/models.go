package api

import (
	"com/gitlab/gituim/repository"
	"encoding/base64"
	"time"
)

type RepositoryModel struct {
	Name   string `json:"name"`
	IsBare bool   `json:"bare"`
}

type RepositoryListModel struct {
	Repositories []string `json:"repositories"`
}

type BranchListModel struct {
	Branches []*BranchModel `json:"branches"`
}

type BranchModel struct {
	Branch string `json:"branch"`
	Commit string `json:"commit"`
	Tree   string `json:"tree"`
	Parent string `json:"parent"`
}

type TagListModel struct {
	Tags []string `json:"tags"`
}

type TagModel struct {
	Tag    string       `json:"tag"`
	Commit *CommitModel `json:"commit"`
}

type CommitModel struct {
	Commit    string          `json:"commit"`
	ShortId   string          `json:"short_id"`
	Tree      string          `json:"tree"`
	Parent    string          `json:"parent"`
	Message   string          `json:"message"`
	Author    *SignatureModel `json:"author"`
	Committer *SignatureModel `json:"committer"`
}

type TreeModel struct {
	Tree    string            `json:"tree"`
	Entries []*TreeEntryModel `json:"entries"`
}

type TreeEntryModel struct {
	Oid      string `json:"oid"`
	FileName string `json:"file_name"`
	Type     string `json:"type"`
}

type SignatureModel struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	When  string `json:"when"`
}

type BlobModel struct {
	Oid      string `json:"oid"`
	ShortId  string `json:"short_id"`
	IsBinary bool   `json:"is_binary"`
	Contents string `json:"contents"`
}

func buildCommitModel(commit *repository.Commit) *CommitModel {
	return &CommitModel{
		Commit:  commit.Commit.String(),
		ShortId: commit.ShortId,
		Parent:  commit.Parent.Id().String(),
		Message: base64.StdEncoding.EncodeToString([]byte(commit.Message)),
		Tree:    commit.Tree.Id().String(),
		Author: &SignatureModel{
			Name:  commit.Author.Name,
			Email: commit.Author.Email,
			When:  commit.Author.When.Format(time.RFC3339),
		},
		Committer: &SignatureModel{
			Name:  commit.Committer.Name,
			Email: commit.Committer.Email,
			When:  commit.Committer.When.Format(time.RFC3339),
		},
	}
}

func buildTreeModel(tree *repository.Tree) *TreeModel {
	var entries []*TreeEntryModel
	for _, entry := range tree.Entries {
		entries = append(entries, &TreeEntryModel{
			Oid:      entry.Id.String(),
			FileName: entry.Name,
			Type:     entry.Type.String(),
		})
	}

	return &TreeModel{
		Tree:    tree.Tree.String(),
		Entries: entries,
	}
}

func buildBlobModel(blob *repository.Blob) *BlobModel {
	return &BlobModel{
		Oid:      blob.Tree.String(),
		ShortId:  blob.ShortId,
		IsBinary: blob.IsBinary,
		Contents: base64.StdEncoding.EncodeToString(blob.Contents),
	}
}

func buildTagModel(tag *repository.Tag) *TagModel {
	return &TagModel{
		Tag:    tag.Tag,
		Commit: buildCommitModel(tag.Commit),
	}
}
