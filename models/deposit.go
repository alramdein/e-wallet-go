package models

import "context"

type DepositRepository interface {
	Create(c context.Context, data CreateDeposit) error
}

type CreateDeposit struct {
	WalletID       int64   `json:"wallet_id"`
	Amount         float64 `json:"amount"`
	AboveThreshold bool    `json:"above_threshold"`
}
