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

	var evt github.Event
	if err = json.Unmarshal(b, &evt); err != nil {
		return nil, err
	}

	return &evt, nil
}
