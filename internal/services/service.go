// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

package services

import (
	"context"
	"errors"

	"coinswallet/internal/domain/wallet/entity"
	"github.com/go-kit/kit/log"
)

type Services interface {

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

type Service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	s := Service{
		logger,
	}
	return s
}

var (
	ErrInService = errors.New("internal service error")

	ErrCreateAccountInvalidName = errors.New("invalid name format")
	ErrCreateAccount            = errors.New("create account error")
	ErrCreateAccountDuplicate   = errors.New("create account error: duplicate name")

	ErrDepositNotFound    = errors.New("account not found")
	ErrDepositAmountError = errors.New("error in amount value")

	ErrTransferFromNotFound = errors.New("from account not found")
	ErrTransferToNotFound   = errors.New("to account not found")
	ErrTransferAmountError  = errors.New("error in amount value")
	ErrTransferNoMoneyError = errors.New("no enough money")

	ErrPaymentsListNotFound         = errors.New("account not found")
	ErrPaymentsListOffsetLimitError = errors.New("error in offset, limit params")

	ErrAccountsListOffsetLimitError = errors.New("error in offset, limit params")
)

func (s Service) CreateAccount(ctx context.Context, name entity.AccountName) (entity.AccountID, error) {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "NewAccount()", "error", err)
		return 0, ErrInService
	}
	if err = a.Validate(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Validate()", "error", err)
		return 0, ErrCreateAccountInvalidName
	}
	if err = a.Register(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Register()", "error", err)
		if a.Find(name) != err { // duplicate
			return 0, ErrCreateAccountDuplicate
		}
		return 0, ErrCreateAccount
	}
	return a.ID, nil
}

func (s Service) Deposit(ctx context.Context, name entity.AccountName, amount float64) (float64, error) {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "NewAccount()", "error", err)
		return 0, ErrInService
	}
	if err = a.Find(name); err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "Find()", "error", err)
		return 0, ErrDepositNotFound
	}
	if amount <= 0 {
		return 0, ErrDepositAmountError
	}
	if _, err = a.Deposit(amount); err != nil {
		_ = s.logger.Log("service", "Deposit", "func", "Deposit()", "error", err)
		return 0, ErrInService
	}
	return a.Balance, nil

}

func (s Service) Transfer(ctx context.Context, from entity.AccountName, to entity.AccountName, amount float64) (entity.ID, error) {
	aFrom, err := entity.NewAccount()
	aTo, _ := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "NewAccount()", "error", err)
		return 0, ErrInService
	}
	if err = aFrom.Find(from); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return 0, ErrTransferFromNotFound
	}
	if err = aTo.Find(to); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return 0, ErrTransferToNotFound
	}
	if amount <= 0 {
		return 0, ErrTransferAmountError
	}

	if aFrom.Balance < amount {
		return 0, ErrTransferNoMoneyError
	}

	var txID int64
	if txID, err = aFrom.Transfer(to, amount); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Transfer()", "error", err)
		return 0, ErrInService
	}
	return entity.ID(txID), nil
}

func (s Service) PaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]entity.Payment, error) {
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

func (s Service) AllPaymentsList(ctx context.Context, offset, limit int64) ([]entity.Payment, error) {
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

func (s Service) AccountsList(ctx context.Context, offset, limit int64) ([]entity.AccountName, error) {
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
