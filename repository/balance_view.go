package repository

import "time"

type BalanceView struct {
	AccountID string `db:"accountID" bson:"accountID" `
	Balance   uint64 `db:"balance" bson:"balance" `
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BalanceViewRepository interface {
	CreateBalanceView(BalanceView) error
	UpsertBulkBalanceView([]BalanceView) error
	UpdateBalanceViewByID(string, BalanceView) error
	GetBalanceViewByID(string) (*BalanceView, error)
	ClearBalanceView() error
}
