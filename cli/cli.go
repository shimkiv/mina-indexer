package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/store"
)

// Run executes the command line interface
func Run() {
	var configPath string
	var runCommand string
	var showVersion bool

	flag.BoolVar(&showVersion, "v", false, "Show application version")
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&runCommand, "cmd", "", "Command to run")
	flag.Parse()

	if showVersion {
		log.Println(config.VersionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		terminate(err)
	}

	if err := initLog(cfg); err != nil {
		terminate(err)
	}

	config.InitRollbar(cfg)
	defer config.TrackRecovery()

	if runCommand == "" {
		terminate("Command is required")
	}

	if err := startCommand(cfg, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, name string) error {
	switch name {
	case "migrate", "migrate:up", "migrate:down", "migrate:redo":
		return startMigrations(name, cfg)
	case "init":
		return startInit(cfg)
	case "server":
		return startServer(cfg)
	case "worker":
		return startWorker(cfg)
	case "sync":
		return runSync(cfg)
	case "status":
		return startStatus(cfg)
	case "update-identity":
		return runUpdateIdentity(cfg)
	default:
		return fmt.Errorf("%s is not a valid command", name)
	}
}

func terminate(message interface{}) {
	if message != nil {
		log.Fatal("ERROR: ", message)
	}
}

func initConfig(path string) (*config.Config, error) {
	cfg := config.New()

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLog(cfg *config.Config) error {
	switch cfg.LogFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		return errors.New("invalid log format: " + cfg.LogFormat)
	}

	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return errors.New("invalid log level: " + cfg.LogLevel)
	}

	return nil
}

func initStore(cfg *config.Config) (*store.Store, error) {
	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	db.SetDebugMode(cfg.LogLevel == "debug")

	return db, nil
}

func initSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return c
}

func initRollbar(cfg *config.Config) {
	config.InitRollbar(cfg)
}
