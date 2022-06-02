package config

import (
	"fmt"
	"strings"

	"github.com/lovoo/goka"
	"github.com/spf13/viper"
)

func GetConf() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}
}

func Env() string {
	return viper.GetString("server.port")
}

func KafkaBrokers() []string {
	return strings.Split(viper.GetString("kafka.brokers"), ",")
}

func KafkaTopic() goka.Stream {
	return goka.Stream(viper.GetString("kafka.topic"))
}

func KafkaGroup() goka.Group {
	return goka.Group(viper.GetString("kafka.group"))
}
