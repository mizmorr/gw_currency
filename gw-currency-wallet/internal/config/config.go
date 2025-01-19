package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Listen

	Storage

	Worker

	Logger

	JWTtokens
}

type Listen struct {
	HttpHost string
	HttpPort string
}

type Storage struct {
	Postgres

	Redis
}

type Postgres struct {
	URL               string
	Timeout           time.Duration
	ConnectAttempts   int
	MaxIdleTime       time.Duration
	MaxOpenConns      int
	HealthCheckPeriod time.Duration
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type Worker struct {
	KeepAliveTimeout time.Duration
}

type Logger struct {
	Level    string
	PathFile string
}

type JWTtokens struct {
	RefreshSecret      string
	AccessSecret       string
	AccessExpiresTime  time.Duration
	RefreshExpiresTime time.Duration
}

var (
	once   sync.Once
	config Config
)

func Get() *Config {
	once.Do(func() {
		viper.AutomaticEnv()

		setDefaults()

		loadConfig()

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})

	return &config
}

func setDefaults() {
	for _, o := range defaults {
		switch o.typing {
		case "string":
			viper.SetDefault(o.name, o.value.(string))
		case "int":
			viper.SetDefault(o.name, o.value.(int))
		default:
			viper.SetDefault(o.name, o.value)
		}
	}
}

func loadConfig() {
	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigName(fileName)
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
}

func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
