package testhelpers

import (
	_ "github.com/lib/pq" // load postgres driver
	"github.com/rs/zerolog/log"
	"os"
)

// IsGithubActions checks that tests are running from github actions.
func IsGithubActions() bool {
	ci := os.Getenv("CI") == "true"
	log.Debug().Msgf("GithubActions: %v", ci)
	return ci
}
