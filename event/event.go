package event

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/alramdein/e-wallet/models"
	"github.com/alramdein/e-wallet/pb"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"google.golang.org/protobuf/encoding/protojson"
)

func InitEmitter(brokers []string, topic goka.Stream) *goka.Emitter {
	emitter, err := goka.NewEmitter(brokers, topic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	return emitter
}

func Send(emitter *goka.Emitter, key string, message string) error {
	// fmt.Println(key, message)
	return emitter.EmitSync(key, message)
	// return nil
}

func RunStreamProcessor(brokers []string, topic goka.Stream, group goka.Group, tmc *goka.TopicManagerConfig, walletUC models.WalletUsecase, wg *sync.WaitGroup) {
	defer wg.Done()
	cb := func(ctx goka.Context, msg interface{}) {
		var counter int64
		if val := ctx.Value(); val != nil {
			counter = val.(int64)
		}
		counter++
		ctx.SetValue(counter)

		log.Printf("key = %s, counter = %v, msg = %v", ctx.Key(), counter, msg)

		pbDeposit := &pb.Deposit{}
		err := protojson.Unmarshal([]byte(msg.(string)), pbDeposit)
		if err != nil {
			log.Fatalf("failed tp decode deposit data: %v", err)
		}

		depositMoney := &models.CreateDeposit{
			WalletID: pbDeposit.WalletId,
			Amount:   float64(pbDeposit.Amount),
		}

		err = walletUC.Deposit(ctx.Context(), *depositMoney)
		if err != nil {
			log.Fatalf("error while deposit the money: %v", err)
		}
	}

	fmt.Println(topic)
	fmt.Println(group)

	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.String), cb),
		goka.Persist(new(codec.Int64)),
	)

	p, err := goka.NewProcessor(brokers,
		g,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	log.Printf("stream process is running...")
	p.Run(context.Background())
}
