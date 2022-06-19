package services

type BalanceView struct {
	AccountID string `json:"accountID" form:"accountID" binding:"required" `
	Balance   uint64 `json:"balance" form:"balance" binding:"required" `
}

type BalanceViewService interface {
	CreateOrUpdateBalanceViewByID(string, BalanceView) error
	GetBalanceViewByID(string) (*BalanceView, error)
	ClearBalanceView() error
}
