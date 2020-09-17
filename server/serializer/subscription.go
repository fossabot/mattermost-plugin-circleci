package serializer

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/chetanyakan/mattermost-plugin-circleci/server/util"
)

type Subscription struct {
	VCSType   string `json:"vcsType"`
	BaseURL   string `json:"baseURL"`
	OrgName   string `json:"orgName"`
	RepoName  string `json:"repoName"`
	ChannelID string `json:"channelID"`
}

// Validate checks if the subscription has valid fields
// returns an error if the subscription is invalid and nil if valid
func (s *Subscription) Validate() error {
	if s.BaseURL != "" {
		if _, err := url.Parse(s.BaseURL); err != nil {
			return errors.Wrap(err, "base url is invalid")
		}
	}

	if s.OrgName == "" {
		return errors.New("org name is empty")
	}

	if s.RepoName == "" {
		return errors.New("repo name is empty")
	}

	return nil
}

// GetKey returns the key against which data can be stored in a map
func (s *Subscription) GetKey() string {
	baseURL, _ := url.Parse(s.BaseURL)
	fields := []string{
		s.VCSType,
		baseURL.Hostname(),
		s.OrgName,
		s.RepoName,
	}
	key := strings.Join(fields, "_")
	return util.GetKeyHash(key)
}
