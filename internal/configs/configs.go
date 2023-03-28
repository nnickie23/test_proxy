package configs

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	setDefaults()

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := unmarshalKeys(config); err != nil {
		return nil, err
	}

	// env vars are in priority
	viper.AutomaticEnv()

	return config, nil
}

type (
	Config struct {
		App      AppConfig
		Postgres PostgresConfig
		Http     HttpConfig
	}

	AppConfig struct {
		Mode     string `mapstructure:"mode"`
		Debug    bool   `mapstructure:"debug"`
		LogLevel string `mapstructure:"log_level"`
	}

	PostgresConfig struct {
		DriverName      string `mapstructure:"driver_name"`
		DataSourceName  string `mapstructure:"data_source_name"`
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}

	HttpConfig struct {
		Addr         string `mapstructure:"addr"`
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
)

func unmarshalKeys(config *Config) error {
	if err := viper.UnmarshalKey("app", &config.App); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &config.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &config.Http); err != nil {
		return err
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("app.mode", defaultAppMode)
	viper.SetDefault("app.debug", defaultAppDebug)
	viper.SetDefault("app.log_level", defaultAppLogLevel)
	viper.SetDefault("postgres.max_open_conns", defaultPostgresMaxOpenConns)
	viper.SetDefault("postgres.max_idle_conns", defaultPostgresMaxIdleConns)
	viper.SetDefault("postgres.conn_max_lifetime", defaultPostgresConnMaxLifetime)
	viper.SetDefault("http.addr", defaultHttpAddr)
	viper.SetDefault("http.read_timeout", defaultHttpReadTimeout)
	viper.SetDefault("http.write_timeout", defaultHttpWriteTimeout)
}

const (
	defaultAppMode     = "dev"
	defaultAppDebug    = false
	defaultAppLogLevel = "warn"

	defaultPostgresMaxOpenConns    = 200
	defaultPostgresMaxIdleConns    = 50
	defaultPostgresConnMaxLifetime = 10 * time.Minute

	defaultHttpAddr         = ":8000"
	defaultHttpReadTimeout  = 5 * time.Minute
	defaultHttpWriteTimeout = 20 * time.Second
)
