package config

import (
	"context"

	"fmt"

	"os"

	"os/signal"

	"syscall"

	"github.com/kelseyhightower/envconfig"
)

const (
	Name = "comptroller"

	defaultIngestHost = "localhost"
	defaultIngestPort = 8080
	ingestAddrFmt     = "%s:%d"

	appIDKey     = "app_id"
	appSecretKey = "app_secret"

	ingestHostKey = "ingest_host"
	ingestPortKey = "ingest_port"
)

type env struct {
	ID     string `default:"not set"`
	Secret string `default:"not set"`

	IngestHost string `envconfig:"ingest_host" default:"localhost"`
	IngestPort int    `envconfig:"ingest_port" default:"8080"`
}

func Init() (context.Context, error) {
	var e env
	if err := envconfig.Process(Name, &e); err != nil {
		return nil, err
	}

	ctx := contextFromEnv(context.Background(), e)
	ctx = contextWithSignalCancel(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	return ctx, nil
}

func contextFromEnv(ctx context.Context, e env) context.Context {
	ctx = context.WithValue(ctx, appIDKey, e.ID)
	ctx = context.WithValue(ctx, appSecretKey, e.Secret)
	ctx = context.WithValue(ctx, ingestHostKey, e.IngestHost)
	ctx = context.WithValue(ctx, ingestPortKey, e.IngestPort)
	return ctx
}

func contextWithSignalCancel(ctx context.Context, signals ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, signals...)

		select {
		case <-sigs:
			cancel()
		case <-ctx.Done():
			// no-op
		}

		signal.Stop(sigs)
		close(sigs)
	}()

	return ctx
}

func IngestAddress(ctx context.Context) string {
	var (
		ok   bool
		host string
		port int
	)

	if host, ok = ctx.Value(ingestHostKey).(string); !ok {
		host = defaultIngestHost
	}

	if port, ok = ctx.Value(ingestPortKey).(int); !ok {
		port = defaultIngestPort
	}

	return fmt.Sprintf(ingestAddrFmt, host, port)
}
