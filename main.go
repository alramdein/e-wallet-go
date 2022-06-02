package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	httpDelivery "github.com/alramdein/e-wallet/delivery/http"
	"github.com/alramdein/e-wallet/event"
	"github.com/alramdein/e-wallet/internal/config"
	"github.com/alramdein/e-wallet/repository"
	"github.com/alramdein/e-wallet/usecase"
	"github.com/labstack/echo/v4"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	brokers []string
	topic   goka.Stream
	group   goka.Group

	tmc *goka.TopicManagerConfig
)

func init() {
	config.GetConf()

	tmc = goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	brokers = config.KafkaBrokers()
	topic = config.KafkaTopic()
	group = config.KafkaGroup()
}

func main() {
	cfg := goka.DefaultConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	goka.ReplaceGlobalConfig(cfg)

	fmt.Println(brokers)
	fmt.Println(config.KafkaGroup())
	fmt.Println(config.KafkaTopic())
	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	err = tm.EnsureStreamExists(string(topic), 8)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", topic, err)
	}

	db, err := leveldb.OpenFile("github.com/alramdein/e-wallet/internal/database/level_db", nil)
	if err != nil {
		log.Fatalf("error instantiating leveldb: %v", err)
	}
	defer db.Close()

	ls, err := storage.New(db)
	if err != nil {
		log.Fatalf("error instantiating leveldb: %v", err)
	}

	e := echo.New()

	walletRepo := repository.NewWalletRepository(ls)
	depositRepo := repository.NewDepositRepository(ls, walletRepo)
	wc := usecase.NewWalletUseacse(walletRepo, depositRepo, time.Duration(2*60))

	emitter := event.InitEmitter(brokers, topic)
	event.RunStreamProcessor(brokers, topic, group, tmc, wc)

	httpDelivery.NewWalletHandler(e, wc, emitter)
	log.Fatal(e.Start(config.Env()))

}
