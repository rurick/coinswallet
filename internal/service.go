package pkg

import (
	"context"

	"coinswallet/pkg/wallet"
)

type Service interface {
	// Send - send amount of currency between two wallet accounts. Return transaction id if success or error
	Send(ctx context.Context, from wallet.AccountName, to wallet.AccountName, amount float64) (wallet.Payment, error)
	// Deposit - deposit amount of currency to the wallet account.  Return transaction id if success or error
	Deposit(ctx context.Context, to wallet.AccountName, amount float64) (wallet.Payment, error)
	// CreateAccount -  create new wallet account
	CreateAccount(ctx context.Context, name wallet.AccountName) error
	// PaymentsList - list of payments of user
	PaymentsList(ctx context.Context, name wallet.AccountName) ([]wallet.Payment, error)
	// AccountsList - List of all registred accounts
	AccountsList(ctx context.Context) ([]wallet.AccountName, error)
}
