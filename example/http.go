package main

import (
	"context"
	"github.com/clearcodecn/lifecycle"
	"github.com/pkg/errors"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	hs := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	lf := lifecycle.New()
	lf.Add(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			ln, _ := net.Listen("tcp", ":3000")
			if err := hs.Serve(ln); err != nil && errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return hs.Shutdown(ctx)
		},
	})
	lf.Start(context.Background())
}
