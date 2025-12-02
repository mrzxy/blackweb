package lib

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var (
	Conf Config
)

type Config struct {
	Port  int    `mapstructure:"port"`
	DSN   string `mapstructure:"dsn"`
	Token string `mapstructure:"token"`
}

var (
	configFile = flag.String("config", "config.yaml", "config file (default is $HOME/.blackweb/config.yaml)")
)

func InitConfig() {
	flag.Parse()
	viper.SetConfigFile(*configFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	Conf = cfg
}
