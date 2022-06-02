package cmd

import (
	"github.com/spf13/cobra"
)

var streamProcessorCmd = &cobra.Command{
	Use:   "stream-processor",
	Short: "run stream-processor",
	Long:  `This subcommand used to run stream-processor`,
	Run:   runStreamProcessor,
}

func init() {
	RootCmd.AddCommand(streamProcessorCmd)
}

/* Change my mind, this is not used anymore. */

func runStreamProcessor(cmd *cobra.Command, args []string) {
	// ls := db.GetLevelDBInstance()
	// walletRepo := repository.NewWalletRepository(ls)
	// depositRepo := repository.NewDepositRepository(ls, walletRepo)
	// walletUC = usecase.NewWalletUseacse(walletRepo, depositRepo, time.Duration(2*60))
	// event.RunStreamProcessor(Brokers, Topic, Group, Tmc, walletUC)
}
