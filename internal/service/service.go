// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

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

	// PaymentsList - list of payments of user.
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all payments
	PaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]entity.Payment, error)

	// AllPaymentsList - list of all payments
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all payments
	AllPaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]entity.Payment, error)

	// AccountsList - List of all registered accounts
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all accounts
	AccountsList(ctx context.Context, offset, limit int64) ([]entity.AccountName, error)
}
