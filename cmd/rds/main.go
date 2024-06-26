package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gwuah/rds/api/gen/proto/v1/protov1connect"
	"github.com/gwuah/rds/internal/cli"
	"github.com/gwuah/rds/internal/config"
	"github.com/gwuah/rds/internal/db"
	"github.com/gwuah/rds/internal/manager"
)

var (
	rootCmd    *cobra.Command
	flagSet    = flag.NewFlagSet("fly-rds", flag.ContinueOnError)
	configFile = flagSet.String("c", "", "Path to config file")
	dbPath     = flagSet.String("db", "rds.db", "Path to sqlite db")
)

func init() {
	flagSet.Parse(os.Args)
	if err := flagSet.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	rootCmd = &cobra.Command{
		Use:              "rds",
		Short:            "remote deployment svc for fly.io",
		Long:             "The remote deployment service for fly.io machine API",
		SilenceUsage:     true,
		TraverseChildren: true,
	}
}

func initLogging() *logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}

func main() {
	args := os.Args
	logger := initLogging()

	if len(args) > 1 && args[1] == "cli" {
		rootCmd.AddCommand(cli.New())
		if err := rootCmd.Execute(); err != nil {
			logger.Fatal(err)
		}
		return
	}

	cfg, err := config.ParseConfigFromFile(*configFile)
	if err != nil {
		logger.Warn("failed to parse config file", err)
		cfg = &config.Config{
			Address: "0.0.0.0",
			Port:    5555,
		}
	}

	db, err := db.New(*dbPath)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup db connection")
	}

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-doneCh
		cancel()
	}()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("okk"))
	}))

	mux.Handle(protov1connect.NewManagerServiceHandler(
		manager.New(logger, db),
	))

	server := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		logger.WithError(err).Fatal("failed to setup tcp listener for server")
	}

	logger.Info("remote.deployment.svc listening on ", listener.Addr())

	go func() {
		server.Serve(listener)
	}()

	<-ctx.Done()

	cancelCtx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	if err := server.Shutdown(cancelCtx); err != nil {
		logger.WithError(err).Error("sremote.deployment.svc shutdown failed")
		return
	}

	logger.Info("remote.deployment.svc shutdown")

}
