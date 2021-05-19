package service

import (
	"context"

	"coinswallet/internal/domain/wallet/entity"
)

type Service interface {
	// Send - send amount of currency between two wallet accounts. Return transaction id if success or error
	Send(ctx context.Context, from entity.AccountName, to entity.AccountName, amount float64) (entity.Payment, error)
	// Deposit - deposit amount of currency to the wallet account.  Return transaction id if success or error
	Deposit(ctx context.Context, to entity.AccountName, amount float64) (entity.Payment, error)
	// CreateAccount -  create new wallet account
	CreateAccount(ctx context.Context, name entity.AccountName) error
	// PaymentsList - list of payments of user
	PaymentsList(ctx context.Context, name entity.AccountName) ([]entity.Payment, error)
	// AccountsList - List of all registred accounts
	AccountsList(ctx context.Context) ([]entity.AccountName, error)
}
