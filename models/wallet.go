package models

import (
	"context"
	"time"
)

type WalletUsecase interface {
	GetWalletByID(c context.Context, id int64) (*Wallet, error)
	Deposit(c context.Context, data CreateDeposit) error
}

type WalletRepository interface {
	GetWalletByID(c context.Context, id int64) (*Wallet, error)
	GetThresholdTimeByIDs(c context.Context, id int64) (*time.Time, error)
}

type Wallet struct {
	ID             int64 `json:"id"`
	Balance        int64 `json:"balance"`
	AboveThreshold bool  `json:"above_threshold"`
}
