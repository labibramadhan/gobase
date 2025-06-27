package start

import (
	"fmt"

	"clodeo.tech/public/go-universe/pkg/config"
	"clodeo.tech/public/go-universe/pkg/env"
	"clodeo.tech/public/go-universe/pkg/logger"
	"github.com/google/gops/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"gobase/di/container"
	"gobase/di/registry"
)

var (
	startCmd = &cobra.Command{
		Use:              "start",
		Short:            "Start all services",
		Long:             "Start all services",
		RunE:             runAllServices,
		PersistentPreRun: rootPreRun,
	}
	serviceName = fmt.Sprintf("%s-%s", "gobase", env.GetEnvironmentName())
)

func rootPreRun(cmd *cobra.Command, args []string) {
	logger.InitGlobalLogger(&logger.Config{
		ServiceName: serviceName,
		Level:       zerolog.DebugLevel,
	})
}

func Cmd() *cobra.Command {
	return startCmd
}

func runAllServices(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")

	app, cleanup, err := container.InitializeApplication(registry.ApplicationContext{
		ConfigPath:  configPath,
		ServiceName: serviceName,
		AgentListen: agent.Listen,
		ReadConfig:  config.ReadConfig,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	app.Run(cleanup)

	return nil
}
