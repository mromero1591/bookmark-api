package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mromero1591/bookmark-api/app/api/config"
	"github.com/mromero1591/bookmark-api/app/api/handlers"
	"github.com/mromero1591/bookmark-api/business/bookmark"
	"github.com/mromero1591/bookmark-api/business/category"
	dbSetup "github.com/mromero1591/bookmark-api/business/database"
	"github.com/mromero1591/bookmark-api/business/sys/metrics"
	"github.com/mromero1591/bookmark-api/business/users"
	"github.com/mromero1591/bookmark-api/database"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/mromero1591/web-foundation/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
func main() {
	// Construct the application logger.
	log, err := logger.New("BOOKMARK-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()
	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}

}

func run(log *zap.SugaredLogger) error {
	// ========================================================================
	// Configuration
	cfg, _ := config.Initialize()

	// ========================================================================
	// App Starting
	expvar.NewString("build").Set(cfg.Build)
	log.Infow("starting service", "version", cfg.Build)
	defer log.Infow("shutdown complete")

	log.Infow("startup", "config", "config")

	// ========================================================================
	// Initialize database
	log.Infow("startup", "status", "initializing database support")
	db, err := dbSetup.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "constructing database")
	}

	queries := database.New(db)

	// ========================================================================
	// Initialize services
	log.Infow("startup", "status", "initializing services")
	userStore := users.NewStore(queries)
	userService := users.NewUserService(userStore)

	categoryStore := category.NewStore(queries)
	categoryService := category.NewCategoryService(categoryStore)

	bookmarkStore := bookmark.NewStore(queries)
	bookmarkService := bookmark.NewBookmarkService(bookmarkStore)

	// ========================================================================
	// Initialize authentication support
	log.Infow("startup", "status", "initializing authentication support")

	a, err := auth.New(cfg.Auth.SigningKey, cfg.Auth.Algorithm)
	if err != nil {
		return errors.Wrap(err, "constructing auth")
	}

	// ===========================================================================
	// Start API Service
	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiConfig := handlers.APIMuxConfig{
		Shutdown:        shutdown,
		Log:             log,
		Metrics:         metrics.New(),
		Auth:            a,
		UserService:     userService,
		CategoryService: categoryService,
		BookmarkService: bookmarkService,
	}
	apiMux := handlers.APIMux(
		apiConfig,
		handlers.WithCORS("OPTIONS"),
	)
	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil

}
