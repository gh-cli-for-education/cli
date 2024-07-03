package ghowner

import (
	"github.com/cli/cli/v2/internal/gh"
)

// Method to get default owner
func GetDefaultOwner(cfg gh.Config) (string, error) {
	optValue := cfg.GetOrDefault("", "gh-owner")
	if optValue.IsSome() {
		return optValue.Unwrap().Value, nil
	}

	return "", nil
}