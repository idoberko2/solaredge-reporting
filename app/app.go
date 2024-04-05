package app

import (
	"context"
	"github.com/idoberko2/semonitor/db"
	"github.com/idoberko2/semonitor/engine"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/seclient"
	"github.com/idoberko2/semonitor/server"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type App interface {
	RunServer(ctx context.Context) error
	RunLastDays(ctx context.Context, rawDays string) error
}

func New() App {
	return &app{}
}

type app struct {
	cfg          general.Config
	engine       engine.Engine
	hcDao        db.HealthCheckDao
	srv          *http.Server
	shutdownDone chan bool
}

func (a *app) RunServer(ctx context.Context) error {
	if err := a.init(); err != nil {
		return err
	}

	if err := a.startServer(ctx); err != nil {
		return err
	}

	return nil
}

func (a *app) RunLastDays(ctx context.Context, rawDays string) error {
	if err := a.init(); err != nil {
		return err
	}

	if err := a.engine.FetchAndPersistLastDays(ctx, rawDays); err != nil {
		return err
	}

	return nil
}

func (a *app) init() error {
	if err := general.LoadDotEnv(); err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}

	cfg, err := general.ReadConfigFromEnv()
	if err != nil {
		return err
	}
	a.cfg = cfg

	migrator := db.NewMigrator()
	if err := migrator.Migrate(cfg); err != nil {
		return err
	}

	eDao := db.NewEnergyDao(cfg)
	if err := eDao.Init(); err != nil {
		return err
	}
	a.hcDao = db.NewHealthCheckDao(cfg)
	if err := a.hcDao.Init(); err != nil {
		return err
	}

	enSvc := engine.NewEnergyService(
		eDao,
		seclient.NewSEClient(req.C(), cfg.SolarEdgeApiKey, cfg.SolarEdgeSiteId),
	)

	eng := engine.New(cfg, enSvc)

	a.engine = eng

	return nil
}

func (a *app) startServer(ctx context.Context) error {
	srv := server.New(a.engine, a.hcDao, a.cfg)
	a.srv = srv

	a.shutdownDone = make(chan bool, 1)
	go a.waitForShutdown(ctx, a.cfg)

	log.Info("starting server...")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// wait for shutdown to complete
	<-a.shutdownDone

	return nil
}

func (a *app) waitForShutdown(ctx context.Context, cfg general.Config) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	sig := <-sigc
	log.
		WithField("signal", sig.String()).
		WithField("timeout", cfg.ServerShutdownTimeout.String()).
		Info("received signal, starting graceful shutdown...")

	ctxWithDeadline, cancel := context.WithTimeout(ctx, cfg.ServerShutdownTimeout)
	defer cancel()

	log.Debug("shutting down...")
	if err := a.srv.Shutdown(ctxWithDeadline); err != nil {
		log.WithError(err).Fatal("Graceful shutdown failed")
	}
	log.Info("shutdown completed successfully")

	a.shutdownDone <- true
}
