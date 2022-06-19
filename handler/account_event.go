package handler

type AccountEventCreateRequest struct {
	AccountID string `json:"accountID" form:"accountID" binding:"required" `
	Money     uint64 `json:"money" form:"money" binding:"required" `
}

type AccountEventTransferRequest struct {
	SenderID   string `json:"senderID" form:"senderID" binding:"required" `
	ReceiverID string `json:"receiverID" form:"receiverID" binding:"required" `
	Money      uint64 `json:"money" form:"money" binding:"required" `
}

type AccountEventHandler interface {
	AccountEventCreate(AccountEventCreateRequest) error
	AccountEventTransfer(AccountEventTransferRequest) error
	ClearAccount() error
}
