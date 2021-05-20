package services

import "coinswallet/internal/domain/wallet/entity"

const (
	PaymentDirectionIncoming = "incoming"
	PaymentDirectionOutgoing = "outgoing"
)

// PaymentEntity using for service response
type PaymentEntity struct {
	Account   entity.AccountName `json:"account"`
	ToAccount entity.AccountName `json:"to_account"`
	Amount    float64            `json:"amount"`
	Direction string             `json:"direction"`
}

// AccountEntity using for service response
type AccountEntity struct {
	Id       entity.AccountName `json:"id"`
	Balance  float64            `json:"balance"`
	Currency string             `json:"currency"`
}
