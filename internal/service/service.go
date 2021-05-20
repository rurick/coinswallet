// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

package service

import (
	"coinswallet/internal/domain/wallet/entity"
	"context"
	"errors"
	"github.com/go-kit/kit/log"
)

type Service interface {

	// CreateAccount -  create new wallet account
	CreateAccount(ctx context.Context, name entity.AccountName) error

	// Deposit - deposit amount of currency to the wallet account.
	Deposit(ctx context.Context, name entity.AccountName, amount float64) error

	// Transfer - send amount of currency between two wallet accounts.
	Transfer(ctx context.Context, from entity.AccountName, to entity.AccountName, amount float64) error

	// PaymentsList - list of payments of the account.
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all payments
	PaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]entity.Payment, error)

	// AllPaymentsList - list of all payments
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all payments
	AllPaymentsList(ctx context.Context, offset, limit int64) ([]entity.Payment, error)

	// AccountsList - List of all registered accounts
	// if set offset and limit > 0 returns slice
	// if limit =-1 returns all accounts
	AccountsList(ctx context.Context, offset, limit int64) ([]entity.AccountName, error)
}

type service struct {
	logger log.Logger
}

func newService(logger log.Logger) service {
	s := service{
		logger,
	}
	return s
}

var (
	ErrInService = errors.New("internal service error")

	ErrCreateAccountInvalidName = errors.New("invalid name format")
	ErrCreateAccount            = errors.New("create account error")

	ErrDepositNotFound    = errors.New("account not found")
	ErrDepositAmountError = errors.New("error in amount value")

	ErrTransferFromNotFound = errors.New("from account not found")
	ErrTransferToNotFound   = errors.New("to account not found")
	ErrTransferAmountError  = errors.New("error in amount value")

	ErrPaymentsListNotFound         = errors.New("account not found")
	ErrPaymentsListOffsetLimitError = errors.New("error in offset, limit params")

	ErrAccountsListOffsetLimitError = errors.New("error in offset, limit params")
)

func (s service) CreateAccount(ctx context.Context, name entity.AccountName) error {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "NewAccount()", "error", err)
		return ErrInService
	}
	if err = a.Validate(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Validate()", "error", err)
		return ErrCreateAccountInvalidName
	}
	if err = a.Register(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Register()", "error", err)
		return ErrCreateAccount
	}
	return nil
}

func (s service) Deposit(ctx context.Context, name entity.AccountName, amount float64) error {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "NewAccount()", "error", err)
		return ErrInService
	}
	if err = a.Find(name); err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "Find()", "error", err)
		return ErrDepositNotFound
	}
	if amount <= 0 {
		return ErrDepositAmountError
	}
	if _, err = a.Deposit(amount); err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "Deposit()", "error", err)
		return ErrInService
	}
	return nil

}

func (s service) Transfer(ctx context.Context, from entity.AccountName, to entity.AccountName, amount float64) error {
	aFrom, err := entity.NewAccount()
	aTo, _ := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "NewAccount()", "error", err)
		return ErrInService
	}
	if err = aFrom.Find(from); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return ErrTransferFromNotFound
	}
	if err = aTo.Find(to); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return ErrTransferToNotFound
	}
	if amount <= 0 {
		return ErrTransferAmountError
	}
	if _, err = aFrom.Transfer(to, amount); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Transfer()", "error", err)
		return ErrInService
	}
	return nil
}

func (s service) PaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]entity.Payment, error) {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "PaymentsList", "func", "NewAccount()", "error", err)
		return nil, ErrInService
	}
	if err = a.Find(name); err != nil {
		_ = s.logger.Log("service", "PaymentsList", "func", "Find()", "error", err)
		return nil, ErrPaymentsListNotFound
	}
	if offset < 0 {
		return nil, ErrPaymentsListOffsetLimitError
	}

	lst, err := entity.PaymentsList(a, offset, limit)
	if err != nil {
		_ = s.logger.Log("service", "PaymentsList", "func", "List()", "error", err)
		return nil, ErrInService
	}

	return lst, nil
}

func (s service) AllPaymentsList(ctx context.Context, offset, limit int64) ([]entity.Payment, error) {
	if offset < 0 {
		return nil, ErrPaymentsListOffsetLimitError
	}

	lst, err := entity.PaymentsList(nil, offset, limit)
	if err != nil {
		_ = s.logger.Log("service", "AllPaymentsList", "func", "List()", "error", err)
		return nil, ErrInService
	}

	return lst, nil
}

func (s service) AccountsList(ctx context.Context, offset, limit int64) ([]entity.AccountName, error) {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "AccountsList", "func", "NewAccount()", "error", err)
		return nil, ErrInService
	}

	if offset < 0 {
		return nil, ErrAccountsListOffsetLimitError
	}

	lst, err := a.List(offset, limit)
	if err != nil {
		_ = s.logger.Log("service", "AccountsList", "func", "List()", "error", err)
		return nil, ErrInService
	}

	return lst, nil
}
