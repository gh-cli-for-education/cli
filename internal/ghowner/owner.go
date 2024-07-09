package ghowner

import (
	"strings"

	"github.com/cli/cli/v2/internal/gh"
)

// TODO: Think a way to get cfg from Factory and not as a parameter
// Method to get default owner
func GetDefaultOwner(cfg gh.Config) (string, error) {
	optValue := cfg.GetOrDefault("", "gh-owner")
	if optValue.IsSome() {
		return optValue.Unwrap().Value, nil
	}

	return "", nil
}

// Method to make REPO to OWNER/REPO if default owner is set
func GetRepoWithDefaultOwner(cfg gh.Config, repo string) (string, error) {
	if repo != "" {
		if !strings.Contains("/", repo) {
			owner, err := GetDefaultOwner(cfg)
			if err != nil {
				return "", err
			}

			if owner != "" {
				return owner + "/" + repo, nil
			}
		}
	}

	return repo, nil
}
