package ghowner

import (
	"strings"
)

// Method to make REPO to OWNER/REPO if default owner is set
func RepoToOwnerRepo(owner string, repo string) (string, error) {
	if repo != "" && !strings.Contains(repo, "/") && owner != "" {
		return owner + "/" + repo, nil
	}

	return repo, nil
}
