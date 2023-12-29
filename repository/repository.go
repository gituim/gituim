package repository

import (
	"fmt"
	"log"
	"os"

	git "github.com/libgit2/git2go/v34"
)

var (
	GCurrentWorkingDirectory = getCurrentWorkingDirectory()
	GRepositoryPrefix        = getEnvOrDefault("GITUIM_REPOSITORY_PREFIX", GCurrentWorkingDirectory)
)

type Repository struct {
	Repository string
	IsBare     bool
}

// ListRepositories - Lists current path repositories, ignores `.git` folders
func ListRepositories() ([]string, error) {
	files, err := os.ReadDir(GRepositoryPrefix)

	if err != nil {
		return nil, fmt.Errorf("unable to list current directory: %w", err)
	}

	var repositories []string

	for _, file := range files {
		if file.IsDir() && file.Name() != ".git" {
			_, err := openRepositoryNoSearch(file.Name())
			if err == nil {
				repositories = append(repositories, file.Name())
			}
		}
	}

	return repositories, nil
}

// CreateRepository - Creates a new bare repository in the current folder
func CreateRepository(repositoryName string) (bool, error) {
	_, err := git.InitRepository(getRepositoryPath(repositoryName), true)

	if err != nil {
		return false, fmt.Errorf("unable to create repository: %w", err)
	}

	log.Printf("Repository %s created", repositoryName)
	return true, nil
}

// DeleteRepository - Delete a repository
func DeleteRepository(repositoryName string) (bool, error) {

	// check if the directory exists first
	if _, err := os.Stat(getRepositoryPath(repositoryName)); os.IsNotExist(err) {
		return false, nil
	}

	err := os.RemoveAll(getRepositoryPath(repositoryName))

	if err != nil {
		return false, fmt.Errorf("unable to delete repository: %w", err)
	}

	log.Printf("Repository %s deleted", repositoryName)
	return true, nil
}

// GetRepositoryInfo - List repository information
func GetRepositoryInfo(repositoryName string) (*Repository, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	return &Repository{
		Repository: repositoryName,
		IsBare:     repository.IsBare(),
	}, nil
}
