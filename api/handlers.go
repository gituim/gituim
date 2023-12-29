package api

import (
	"com/gitlab/gituim/repository"
	"encoding/json"
	"io"
	"net/http"
)

func ListRepositoriesHandler(w http.ResponseWriter, _ *http.Request) {
	repos, err := repository.ListRepositories()
	if err != nil {
		handleError(err, w)
		return
	}

	if repos != nil {
		data, err := json.Marshal(RepositoryListModel{Repositories: repos})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func CreateRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to parse repository name", http.StatusInternalServerError)
		return
	}

	var repo RepositoryModel
	err = json.Unmarshal(body, &repo)
	if err != nil || repo.Name == "" {
		http.Error(w, "invalid repository name", http.StatusBadRequest)
		return
	}

	_, err = repository.CreateRepository(repo.Name)
	if err != nil {
		handleError(err, w)
		return
	}

	w.Header().Add("Location", repo.Name)
	w.WriteHeader(http.StatusOK)
}

func GetRepositoryInfoHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	info, err := repository.GetRepositoryInfo(repositoryName)
	if err != nil {
		handleError(err, w)
		return
	}

	data, err := json.Marshal(RepositoryModel{
		Name:   info.Repository,
		IsBare: info.IsBare})
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		handleError(err, w)
		return
	}
}

func DeleteRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	deleted, err := repository.DeleteRepository(repositoryName)
	if err != nil {
		handleError(err, w)
		return
	}

	if deleted {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func ListBranchesHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	branches, err := repository.ListRepositoryBranches(repositoryName)
	if err != nil {
		handleError(err, w)
		return
	}

	if branches != nil {
		dto := BranchListModel{}
		for _, branch := range branches {
			dto.Branches = append(dto.Branches, &BranchModel{
				Branch: branch.Branch,
				Commit: branch.Commit.String(),
				Tree:   branch.Tree.String(),
				Parent: branch.Parent.String(),
			})
		}

		data, err := json.Marshal(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetBranchHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	branchName, ok := getVar(w, r, "branch")
	if !ok {
		return
	}

	branch, err := repository.GetBranch(repositoryName, branchName)
	if err != nil {
		handleError(err, w)
		return
	}

	if branch != nil {
		data, err := json.Marshal(BranchModel{
			Branch: branchName,
			Commit: branch.Commit.String(),
			Tree:   branch.Tree.String(),
			Parent: branch.Parent.String()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetCommitHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	commitOid, ok := getVar(w, r, "commit")
	if !ok {
		return
	}

	commit, err := repository.LookupCommit(repositoryName, commitOid)
	if err != nil {
		handleError(err, w)
		return
	}

	data, err := json.Marshal(buildCommitModel(commit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTreeHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	treeOid, ok := getVar(w, r, "tree")
	if !ok {
		return
	}

	tree, err := repository.LookupTree(repositoryName, treeOid)
	if err != nil {
		handleError(err, w)
		return
	}

	data, err := json.Marshal(buildTreeModel(tree))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetBlobHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	blobOid, ok := getVar(w, r, "blob")
	if !ok {
		return
	}

	blob, err := repository.LookupBlob(repositoryName, blobOid)
	if err != nil {
		handleError(err, w)
		return
	}

	data, err := json.Marshal(buildBlobModel(blob))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	tags, err := repository.ListRepositoryTags(repositoryName)
	if err != nil {
		handleError(err, w)
		return
	}

	if len(tags) > 0 {
		data, err := json.Marshal(TagListModel{Tags: tags})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetTagHandler(w http.ResponseWriter, r *http.Request) {
	repositoryName, ok := getVar(w, r, "repository")
	if !ok {
		return
	}

	tagName, ok := getVar(w, r, "tag")
	if !ok {
		return
	}

	tag, err := repository.LookupTag(repositoryName, tagName)
	if err != nil {
		handleError(err, w)
		return
	}

	data, err := json.Marshal(buildTagModel(tag))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
