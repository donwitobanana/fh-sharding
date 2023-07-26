package config

import (
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	WorkersMaxNumber        int `env:"WORKERS_MAX_NUMBER"`
	WorkersInitNumber       int `env:"WORKERS_INIT_NUMBER"`
	WorkersTimeoutInSeconds int `env:"WORKERS_TIMEOUT_IN_SECONDS"`
	WorkersBatchSize        int `env:"WORKERS_BATCH_SIZE"`
	PartitionsNumber        int `env:"PARTITIONS_NUMBER"`
}

func LoadAppConfig() (AppConfig, error) {
	var config AppConfig

	opts := env.Options{RequiredIfNoDef: true}
	err := env.ParseWithOptions(&config, opts)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
