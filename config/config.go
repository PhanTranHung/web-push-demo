package config

import (
	"fmt"
)

type Configs struct {
	env EnvConfigs
}

type ConfigLoader func(config *Configs) error

func LoadConfig(loaders ...ConfigLoader) Configs {

	c := Configs{}

	// Use all loader to load environment variable to Setting
	for _, load := range loaders {
		if err := load(&c); err != nil {
			panic(fmt.Errorf("failed to load configs, error : %w", err))
		}
	}
	return c
}

func (c *Configs) GetENVConfigs() EnvConfigs {
	return c.env
}
