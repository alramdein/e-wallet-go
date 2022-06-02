package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alramdein/e-wallet/models"
	"github.com/lovoo/goka/storage"
)

type depositRepo struct {
	db         storage.Storage
	walletRepo models.WalletRepository
}

func NewDepositRepository(db storage.Storage, walletRepo models.WalletRepository) models.DepositRepository {
	return &depositRepo{
		db:         db,
		walletRepo: walletRepo,
	}
}

func (d *depositRepo) Create(c context.Context, data models.CreateDeposit) error {
	wallet, err := d.walletRepo.GetWalletByID(c, data.WalletID)
	if err != nil {
		return err
	}

	wallet.Balance += int64(data.Amount)
	mwallet, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	err = d.db.Set(fmt.Sprint(data.WalletID), mwallet)
	if err != nil {
		return err
	}

	return nil
}
