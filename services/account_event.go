package services

type AccountEventCreateRequest struct {
	AccountID string `json:"accountID" form:"accountID" binding:"required" `
	Money     uint64 `json:"money" form:"money" binding:"required" `
}

// uint8
type AccountEventTransferRequest struct {
	SenderID   string `json:"senderID" form:"senderID" binding:"required" `
	ReceiverID string `json:"receiverID" form:"receiverID" binding:"required" `
	Money      uint64 `json:"money" form:"money" binding:"required" `
}

type AccountEventService interface {
	NewAccountEvent(AccountEventCreateRequest) error
	AccountEventTransfer(AccountEventTransferRequest) error
	ClearAccount() error
}
