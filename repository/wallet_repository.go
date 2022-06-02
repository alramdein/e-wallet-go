package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	walletByte, err := w.db.Get(fmt.Sprint(id))
	if err != nil {
		return nil, err
	}

	wallet := &models.Wallet{}
	err = json.Unmarshal(walletByte, &wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

/* get the time when deposit amount surpass the threshold */
func (w *walletRepo) GetThresholdTimeByIDs(c context.Context, id int64) (*time.Time, error) {
	ttByte, err := w.db.Get(w.thresholdKey(id))
	if err != nil {
		return nil, err
	}

	tt := &time.Time{}
	err = json.Unmarshal(ttByte, &tt)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (w *walletRepo) thresholdKey(id int64) string {
	return fmt.Sprintf("threshold:wallet:%d", id)
}
