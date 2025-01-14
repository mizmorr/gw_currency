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
	LogLevel string

	Listen

	Storage
}

type Listen struct {
	Host string
	Port int
}

type Storage struct {
	PostgresURL               string
	PostgresTimeout           time.Duration
	PostgresConnectAttempts   int
	PostgresMaxIdleTime       time.Duration
	PostgresMaxOpenConns      int
	PostgresHealthCheckPeriod time.Duration
}

var (
	once   sync.Once
	config Config
)

func Get() *Config {
	once.Do(func() {
		viper.AutomaticEnv()

		for _, o := range Defaults {
			switch o.Typing {
			case "string":
				fmt.Println(o.Name)
				viper.SetDefault(o.Name, o.Value.(string))
			case "int":
				viper.SetDefault(o.Name, o.Value.(int))
			default:
				viper.SetDefault(o.Name, o.Value)
			}
		}
		if fileName := viper.GetString("config"); fileName != "" {
			viper.SetConfigName(fileName)
			viper.SetConfigType("env")
			viper.AddConfigPath(".")

			if err := viper.ReadInConfig(); err != nil {
				panic(err)
			}
		}
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})

	return &config
}

func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
