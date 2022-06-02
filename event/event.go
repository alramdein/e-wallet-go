package event

import (
	"context"
	"log"

	"github.com/alramdein/e-wallet/models"
	"github.com/alramdein/e-wallet/pb"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"google.golang.org/protobuf/proto"
)

func InitEmitter(brokers []string, topic goka.Stream) *goka.Emitter {
	emitter, err := goka.NewEmitter(brokers, topic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitter.Finish()
	return emitter
}

func Send(emitter *goka.Emitter, key string, message []byte) error {
	return emitter.EmitSync(key, message)
}

func RunStreamProcessor(brokers []string, topic goka.Stream, group goka.Group, tmc *goka.TopicManagerConfig, walletUC models.WalletUsecase) {
	cb := func(ctx goka.Context, msg interface{}) {
		var counter int64
		if val := ctx.Value(); val != nil {
			counter = val.(int64)
		}
		counter++
		ctx.SetValue(counter)

		log.Printf("key = %s, counter = %v, msg = %v", ctx.Key(), counter, msg)

		pbDeposit := &pb.Deposit{}
		depositMoney := models.CreateDeposit{}
		err := proto.Unmarshal(msg.([]byte), pbDeposit)
		if err != nil {
			log.Fatalf("failed tp decode deposit data: %v", err)
		}
		err = walletUC.Deposit(ctx.Context(), depositMoney)
		if err != nil {
			log.Fatalf("error while deposit the money: %v", err)
		}
	}

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
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Printf("error running processor: %v", err)
		}
	}()

	<-done
	cancel()
	<-done
}
