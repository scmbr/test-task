package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Hasher   HasherConfig
		Auth     AuthConfig
	}
	HTTPConfig struct {
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
	PostgresConfig struct {
		Username string `mapstructure:"username"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		Password string
	}
	HasherConfig struct {
		Cost int `mapstructure:"cost"`
	}
	AuthConfig struct {
		SigningKey      string
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	}
)

func Init(configsDir string) (*Config, error) {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}
func unmarshal(cfg *Config) error {

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("hasher", &cfg.Hasher); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	return nil
}
func setFromEnv(cfg *Config) {
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Auth.SigningKey = os.Getenv("SIGNING_KEY")
}
