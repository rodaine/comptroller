package gh

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/google/go-github/github"
)

func Extract(req *http.Request) (*github.Event, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	raw := struct {
		Repo *github.Repository `json:"repository,omitempty"`
	}{}

	if err = json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	typ := github.WebHookType(req)
	return &github.Event{
		RawPayload: (*json.RawMessage)(&b),
		Repo:       raw.Repo,
		Type:       &typ,
	}, nil
}
