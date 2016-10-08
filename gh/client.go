package gh

import (
	"log"

	"context"

	"github.com/google/go-github/github"
)

type Client interface {
	Zen() (msg string, resp *github.Response, err error)
}

type client struct {
	*github.Client
}

func Init(ctx context.Context) (Client, error) {
	c := github.NewClient(nil)

	if msg, _, err := c.Zen(); err != nil {
		return nil, err
	} else {
		log.Println(msg)
	}

	return &client{Client: c}, nil
}
