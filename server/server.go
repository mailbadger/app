package server

import (
	"context"
	"net/http"
	"time"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/routes"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Addr    string
	Key     string
	Cert    string
	Handler http.Handler
}

func From(api routes.API, conf config.Config) *Server {
	return New(
		":"+conf.Server.Port,
		api.Handler(),
		conf.Server.Cert,
		conf.Server.Key,
	)
}

func New(addr string, handler http.Handler, cert, key string) *Server {
	return &Server{
		Addr:    addr,
		Cert:    cert,
		Key:     key,
		Handler: handler,
	}
}

func (s Server) ListenAndServe(ctx context.Context) error {
	if s.Key != "" {
		return s.listenAndServeTLS(ctx)
	}
	err := s.listenAndServe(ctx)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s Server) listenAndServe(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}
	g.Go(func() error {
		<-ctx.Done()

		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		return s1.Shutdown(ctxShutdown)
	})
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	return g.Wait()
}

func (s Server) listenAndServeTLS(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(redirect),
	}
	s2 := &http.Server{
		Addr:    ":https",
		Handler: s.Handler,
	}
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	g.Go(func() error {
		return s2.ListenAndServeTLS(
			s.Cert,
			s.Key,
		)
	})
	g.Go(func() error {
		<-ctx.Done()

		var gShutdown errgroup.Group
		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		gShutdown.Go(func() error {
			return s1.Shutdown(ctxShutdown)
		})
		gShutdown.Go(func() error {
			return s2.Shutdown(ctxShutdown)
		})

		return gShutdown.Wait()
	})
	return g.Wait()
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
