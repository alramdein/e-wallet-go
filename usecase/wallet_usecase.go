package usecase

import (
	"context"
	"time"

	"github.com/alramdein/e-wallet/models"
)

type walletUsecase struct {
	walletRepo     models.WalletRepository
	depositRepo    models.DepositRepository
	contextTimeout time.Duration
}

func NewWalletUseacse(walletRepo models.WalletRepository, depositRepo models.DepositRepository, timeout time.Duration) models.WalletUsecase {
	return &walletUsecase{
		walletRepo:     walletRepo,
		depositRepo:    depositRepo,
		contextTimeout: timeout,
	}
}

func (w *walletUsecase) GetWalletByID(c context.Context, id int64) (*models.Wallet, error) {
	wallet, err := w.walletRepo.GetWalletByID(c, id)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrNotFound
	}
	return wallet, nil
}

func (w *walletUsecase) Deposit(c context.Context, data models.CreateDeposit) error {
	ctx, cancel := context.WithTimeout(c, w.contextTimeout)
	defer cancel()
	_, err := w.GetWalletByID(ctx, data.WalletID)
	if err != nil {
		return err
	}
	err = w.depositRepo.Create(ctx, models.CreateDeposit{
		WalletID: data.WalletID,
		Amount:   data.Amount,
	})
	if err != nil {
		return err
	}
	return nil
}
