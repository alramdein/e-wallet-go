/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/alramdein/e-wallet/internal/config"
	"github.com/lovoo/goka"
	"github.com/spf13/cobra"
)

var (
	Brokers []string
	Topic   goka.Stream
	Group   goka.Group

	Tmc *goka.TopicManagerConfig
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "e-wallet",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.GetConf()

	Tmc = goka.NewTopicManagerConfig()
	Tmc.Table.Replication = 1
	Tmc.Stream.Replication = 1

	Brokers = config.KafkaBrokers()
	Topic = config.KafkaTopic()
	Group = config.KafkaGroup()

	cfg := goka.DefaultConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	goka.ReplaceGlobalConfig(cfg)

	Tm, err := goka.NewTopicManager(Brokers, goka.DefaultConfig(), Tmc)

	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	err = Tm.EnsureStreamExists(string(Topic), 8)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", Topic, err)
	}
}
