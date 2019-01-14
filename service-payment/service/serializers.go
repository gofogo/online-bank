package service

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentEvent struct {
	PayementUId           string    `json:"paymentUId"`
	DestinationAccountUid string    `json:"destinationAccountUId"`
	PayeeUId              string    `json:"payeeUId"`
	Amount                int64     `json:"amount"`
	Created               time.Time `json:"created"`
	Status                string    `json:"status"`
}

type CreatePaymentFeed struct {
	DestinationAccountUid string `json:"destinationAccountUId" binding:"required,uuid4"`
	PayeeUId              string `json:"payeeUId" binding:"required,uuid4"`
	Amount                int64  `json:"amount" binding:"required,gt=0"`
}

// TODO: Q>why cannot user *PaymentRegister?
func Event(p CreatePaymentFeed) CreatePaymentEvent {
	paymentuid := uuid.New() // Utils should generate that
	result := CreatePaymentEvent{
		PayementUId:           paymentuid.String(),
		DestinationAccountUid: p.DestinationAccountUid,
		PayeeUId:              p.PayeeUId,
		Amount:                p.Amount,
		Created:               time.Now(),
		Status:                "accepted",
	}
	return result
}
