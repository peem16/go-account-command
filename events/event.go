package events

import "reflect"

var Topics = []string{
	reflect.TypeOf(AccountCreateEvent{}).Name(),
	reflect.TypeOf(AccountTransferEvent{}).Name(),
}

type Event interface {
}

type AccountCreateEvent struct {
	ID        string
	AccountID string
	Money     uint64
	Type      string
}

type AccountTransferEvent struct {
	ID         string
	SenderID   string
	ReceiverID string
	Money      uint64
	Type       string
}
