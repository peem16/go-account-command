package services

import (
	"encoding/json"
	"go-account-command/errs"
	"go-account-command/events"
	"go-account-command/logs"
	"go-account-command/messageBroker"
	"go-account-command/repository"

	"log"
	"reflect"
)

type balanceViewService struct {
	balanceViewRepo  repository.BalanceViewRepository
	accountEventRepo repository.AccountEventRepository
}

func NewBalanceViewService(balanceViewRepo repository.BalanceViewRepository, accountEventRepo repository.AccountEventRepository) messageBroker.EventHandler {
	return balanceViewService{balanceViewRepo: balanceViewRepo, accountEventRepo: accountEventRepo}
}

func (b balanceViewService) CreateOrUpdateBalanceViewByID(id string, a BalanceView) error {
	repoBalanceView := repository.BalanceView{
		AccountID: a.AccountID,
		Balance:   a.Balance,
	}

	err := b.balanceViewRepo.CreateBalanceView(repoBalanceView)

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}

	return nil
}

func (b balanceViewService) GetBalanceViewByID(string) (*BalanceView, error) {

	return nil, nil
}

func (b balanceViewService) ClearBalanceView() error {

	err := b.balanceViewRepo.ClearBalanceView()

	if err != nil {
		logs.Error(err)
		return errs.NewUnexpectedError()
	}
	return nil
}

func (obj balanceViewService) Handle(topic string, eventBytes []byte) {
	switch topic {
	case reflect.TypeOf(events.AccountCreateEvent{}).Name():
		event := &events.AccountCreateEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}
		balanceView := repository.BalanceView{
			AccountID: event.AccountID,
			Balance:   uint64(event.Money),
		}
		err = obj.balanceViewRepo.CreateBalanceView(balanceView)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %#v", topic, event)
	case reflect.TypeOf(events.AccountTransferEvent{}).Name():
		event := &events.AccountTransferEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			log.Println(err)
			return
		}

		arr := []string{event.SenderID, event.ReceiverID}

		accountEventData, err := obj.accountEventRepo.GetAccountEvents(arr)

		if err != nil {
			log.Println(err)
		}

		mapBalanceView := map[string]BalanceView{}

		for _, evet := range accountEventData {

			if evet.Type == "CREATE" {

				mapBalanceView[evet.AccountID] = BalanceView{evet.AccountID, uint64(evet.Money)}
			} else if evet.Type == "TRANSFER" {

				mapBalanceView[evet.AccountID] = BalanceView{evet.AccountID, uint64(mapBalanceView[evet.AccountID].Balance - uint64(evet.Money))}
				mapBalanceView[evet.ReceiverID] = BalanceView{evet.ReceiverID, uint64(mapBalanceView[evet.ReceiverID].Balance + uint64(evet.Money))}
			}
		}

		balanceViews := []repository.BalanceView{
			{AccountID: mapBalanceView[event.SenderID].AccountID, Balance: uint64(mapBalanceView[event.SenderID].Balance)},
			{AccountID: mapBalanceView[event.ReceiverID].AccountID, Balance: uint64(mapBalanceView[event.ReceiverID].Balance)},
		}

		err = obj.balanceViewRepo.UpsertBulkBalanceView(balanceViews)

		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%v] %#v", topic, event)

	default:
		log.Println("no event handler")
	}
}
