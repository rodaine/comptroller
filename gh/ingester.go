package gh

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/rodaine/comptroller/config"
	"github.com/tylerb/graceful"
)

func Ingest(ctx context.Context) (<-chan *github.Event, <-chan error) {
	events := make(chan *github.Event)
	errors := make(chan error)

	in := &Ingester{
		events: events,
		errors: errors,
	}

	go func() {
		if err := in.Listen(ctx); err != nil {
			errors <- err
		}

		close(events)
		close(errors)
	}()

	return events, errors
}

type Ingester struct {
	s *graceful.Server

	events chan<- *github.Event
	errors chan<- error
}

func (in *Ingester) Listen(ctx context.Context) error {
	in.s = &graceful.Server{
		Server: &http.Server{
			Addr:    config.IngestAddress(ctx),
			Handler: in.handler(),
		},
		NoSignalHandling: true,
	}

	go func() {
		<-ctx.Done()
		in.s.Stop(5 * time.Second)
	}()

	return in.s.ListenAndServe()
}

func (in *Ingester) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ingest/", in.handleHook)
	return mux
}

func (in *Ingester) handleHook(res http.ResponseWriter, req *http.Request) {
	defer res.WriteHeader(http.StatusNoContent)
	defer req.Body.Close()

	evt, err := Extract(req)
	if err != nil {
		in.errors <- err
		return
	}

	in.events <- evt
}
