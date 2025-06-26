package provider

import (
	"github.com/google/gops/agent"
	"github.com/rs/zerolog/log"

	"gobase/config"
	"gobase/di/registry"
)

var (
	errInitConfig = "failed to initiate config"
)

func ProvideConfig(
	appContext registry.ApplicationContext,
	agentListen registry.AgentListen,
	readConfig registry.ReadConfig,
) *config.MainConfig {
	if err := agentListen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal().Err(err).Msg(errInitConfig)
	}

	cfg := &config.MainConfig{}
	log.Info().Msgf("reading config from %s", appContext.ConfigPath)
	err := readConfig(cfg, appContext.ConfigPath, "config")
	if err != nil {
		log.Fatal().Err(err).Msg(errInitConfig)
	}

	return cfg

}
