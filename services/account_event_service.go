package services

import (
	"fmt"
	"go-account-command/errs"
	"go-account-command/events"
	"go-account-command/logs"
	"go-account-command/messageBroker"
	"go-account-command/repository"
	"time"

	"github.com/google/uuid"
)

type accountEventService struct {
	accountEventRepo repository.AccountEventRepository
	balanceEventRepo repository.BalanceViewRepository
	eventProducer    messageBroker.EventProducer
}

func NewAccountEventService(accountEventRepo repository.AccountEventRepository, balanceEventRepo repository.BalanceViewRepository, eventProducer messageBroker.EventProducer) AccountEventService {
	return accountEventService{accountEventRepo: accountEventRepo, balanceEventRepo: balanceEventRepo, eventProducer: eventProducer}
}

func (obj accountEventService) NewAccountEvent(a AccountEventCreateRequest) error {

	types := []string{"CREATE"}
	accountEventData, err := obj.accountEventRepo.FindOneAccountEventByID(a.AccountID, types)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	fmt.Println(accountEventData)
	if accountEventData != nil {
		return errs.NewValidationError("Account does exist.")
	}

	repoAccountEvent := repository.AccountEvent{
		AccountID: a.AccountID,
		Money:     a.Money,
		Type:      "CREATE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = obj.accountEventRepo.CreateEvent(repoAccountEvent)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	event := events.AccountCreateEvent{
		ID:        uuid.NewString(),
		AccountID: a.AccountID,
		Money:     a.Money,
		Type:      "CREATE",
	}

	err = obj.eventProducer.Produce(event)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (obj accountEventService) AccountEventTransfer(a AccountEventTransferRequest) error {

	if a.SenderID == a.ReceiverID {
		return errs.NewValidationError("can't send to yourself")
	}

	balanceSender, err := obj.balanceEventRepo.GetBalanceViewByID(a.SenderID)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	if balanceSender == nil {
		return errs.NewValidationError("senderID does not exist.")
	}

	if int(balanceSender.Balance-a.Money) < 0 {
		return errs.NewValidationError("Your balance is not enough.")
	}

	balanceReceiver, err := obj.balanceEventRepo.GetBalanceViewByID(a.ReceiverID)
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	if balanceReceiver == nil {
		return errs.NewValidationError("receiverID does not exist.")
	}

	repoAccountEvent := repository.AccountEvent{
		AccountID:  a.SenderID,
		ReceiverID: a.ReceiverID,
		Money:      a.Money,
		Type:       "TRANSFER",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = obj.accountEventRepo.CreateEvent(repoAccountEvent)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	event := events.AccountTransferEvent{
		ID:         uuid.NewString(),
		SenderID:   a.SenderID,
		ReceiverID: a.ReceiverID,
		Money:      a.Money,
		Type:       "TRANSFER",
	}

	err = obj.eventProducer.Produce(event)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (t accountEventService) ClearAccount() error {

	err := t.accountEventRepo.ClearAccount()

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	err = t.balanceEventRepo.ClearBalanceView()
	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}
