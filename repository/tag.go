package repository

// ListRepositoryTags - list all tags in the repository
func ListRepositoryTags(repositoryName string) ([]string, error) {
	repository, err := openRepositoryNoSearch(repositoryName)
	if err != nil {
		return nil, handleGitError(err, "unable to open repository")
	}

	tags, err := repository.Tags.List()
	if err != nil {
		return nil, handleGitError(err, "unable to create branch iterator")
	}

	return tags, nil
}
