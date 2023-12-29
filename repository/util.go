package repository

import (
	"fmt"
	"os"

	git "github.com/libgit2/git2go/v34"
)

// Open a repository with Bare and NoSearch flags enabled
func openRepositoryNoSearch(repositoryName string) (*git.Repository, error) {
	flags := git.RepositoryOpenBare | git.RepositoryOpenNoSearch
	return git.OpenRepositoryExtended(getRepositoryPath(repositoryName), flags, "")
}

func getEnvOrDefault(key, fallbackValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallbackValue
}

func getCurrentWorkingDirectory() string {
	if dir, err := os.Getwd(); err != nil {
		panic("unable to get working directory")
	} else {
		return dir
	}
}

func getRepositoryPath(repositoryName string) string {
	return fmt.Sprintf("%s/%s", GRepositoryPrefix, repositoryName)
}
