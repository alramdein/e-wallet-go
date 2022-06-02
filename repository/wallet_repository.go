package repository

import (
	"context"
	"encoding/json"

	"github.com/alramdein/e-wallet/models"
	"github.com/lovoo/goka/storage"
)

type walletRepo struct {
	db storage.Storage
}

func NewWalletRepository(db storage.Storage) models.WalletRepository {
	return &walletRepo{
		db: db,
	}
}

func (w *walletRepo) GetWalletByID(c context.Context, id int64) (*models.Wallet, error) {
	walletByte, err := w.db.Get(string(id))
	if err != nil {
		return nil, err
	}
	wallet := &models.Wallet{}
	json.Unmarshal(walletByte, wallet)
	return wallet, nil
}
