package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"leggett.dev/devmarks/api/api"
	"leggett.dev/devmarks/api/app"
)

type query struct{}

func serveAPI(ctx context.Context, api *api.API) {
	router := mux.NewRouter()
	api.Init(router.PathPrefix("/api").Subrouter())
	api.InitGraphql(router)

	var server *http.Server
	var handler http.Handler

	if api.Config.Cors {
		cors := handlers.CORS(
		handlers.AllowedOrigins(api.Config.AllowedHosts),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	handler = cors(router)
	
} else {
	handler = router
}
server = &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Config.Port),
		Handler:     handler,
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}
		close(done)
	}()

	logrus.Infof("serving api at http://127.0.0.1:%d", api.Config.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Error(err)
	}
	<-done
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the api",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New()
		if err != nil {
			return err
		}
		defer a.Close()

		api, err := api.New(a)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info("signal caught. shutting down...")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveAPI(ctx, api)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
