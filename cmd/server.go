package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	httpDelivery "github.com/alramdein/e-wallet/delivery/http"
	"github.com/alramdein/e-wallet/event"
	"github.com/alramdein/e-wallet/internal/config"
	db "github.com/alramdein/e-wallet/internal/database"
	"github.com/alramdein/e-wallet/models"
	"github.com/alramdein/e-wallet/repository"
	"github.com/alramdein/e-wallet/usecase"
	"github.com/labstack/echo/v4"
	"github.com/lovoo/goka/storage"
	"github.com/spf13/cobra"
)

var (
	walletUC models.WalletUsecase
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup

	db.InitLevelDB()
	ls := db.GetLevelDBInstance()
	initTestwallet(ls)

	e := echo.New()

	walletRepo := repository.NewWalletRepository(ls)
	depositRepo := repository.NewDepositRepository(ls, walletRepo)
	walletUC = usecase.NewWalletUseacse(walletRepo, depositRepo, time.Duration(2*60))

	emitter := event.InitEmitter(Brokers, Topic)

	httpDelivery.NewWalletHandler(e, walletUC, emitter)

	wg.Add(2)
	go startWebServer(e, &wg)
	go event.RunStreamProcessor(Brokers, Topic, Group, Tmc, walletUC, &wg)
	wg.Wait()
}

func startWebServer(e *echo.Echo, wg *sync.WaitGroup) {
	defer wg.Wait()
	e.Start(config.Env())
}

/* init wallet for test purpose */
func initTestwallet(ls storage.Storage) {
	wallet := &models.Wallet{
		ID:             1,
		Balance:        0,
		AboveThreshold: false,
	}
	w, err := json.Marshal(wallet)
	if err != nil {
		log.Fatalf("failed to decode wallet for test")
	}

	ls.Set(fmt.Sprint(wallet.ID), w)
}
