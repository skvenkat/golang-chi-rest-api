package app

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// This prefix is used for environment variables to override config file values
const envVarPrefix = "APISERVER"

const baseDir = "./configs"

func LoadConfig(deployment string) *Config {
	var cfg *Config = nil
	opt := loadOptions{
		EnvVarPrefix: envVarPrefix,
		Deployment:   deployment,
		BaseDir:      baseDir,
	}
	cfg = &Config{}
	err := loadFromYaml(opt, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration file: %v", err))
	}
	cfg.Deployment = deployment
	return cfg
}

type loadOptions struct {
	EnvVarPrefix string
	Deployment   string
	BaseDir      string
}

func loadFromYaml(opt loadOptions, cfg *Config) error {
	v := viper.New()
	v.SetEnvPrefix(opt.EnvVarPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName(opt.Deployment)
	v.SetConfigType("yaml")
	v.AddConfigPath(opt.BaseDir)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading configuration file: %w", err)
	}
	err = v.UnmarshalExact(cfg)
	if err != nil {
		return fmt.Errorf("error parsing configuration file: %w", err)
	}
	return nil
}
