package main

import (
	"flag"
	"fmt"
	"os"

	"app/agent_grid/cmd/app"
	"app/agent_grid/internal/config"

	"bitbucket.org/fyscal/be-commons/pkg/log"
)

var Version = ""
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	fmt.Printf("Loading application configuration...\n")
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		_ = fmt.Errorf("Failed to load application configuration: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("Application configuration loaded successfully.")

	// Initialize logger before loading the configuration
	var logger log.Logger = log.New(
		log.LogConfig{
			ServiceName: cfg.AppName,
			AppEnv:      cfg.Environment,
			AppVersion:  fmt.Sprintf("%s(%s)", cfg.AppVersion, Version),
		},
	).With(nil)

	if logger == nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to initialize logger.\n")
		os.Exit(-1)
	}

	// Initialize and start the application
	application := &app.App{}
	application.NewApplication(cfg, logger)
	fmt.Printf("Application initialized successfully.")
}
