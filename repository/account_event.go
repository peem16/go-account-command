package repository

import "time"

type AccountEvent struct {
	AccountID  string `db:"accountID" bson:"accountID" `
	ReceiverID string `db:"receiverID" bson:"receiverID,omitempty"`
	Money      uint64 `db:"money" `
	Type       string `db:"type" `
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AccountEventRepository interface {
	CreateEvent(AccountEvent) error
	GetAccountEventByID(string) ([]AccountEvent, error)
	GetAccountEvents([]string) ([]AccountEvent, error)
	FindOneAccountEventByID(string, []string) (*AccountEvent, error)
	ClearAccount() error
}
