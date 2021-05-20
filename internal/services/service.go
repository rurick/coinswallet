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
	AccountsList(ctx context.Context, offset, limit int64) ([]entity.Account, error)
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

	ErrTransferFromNotFound    = errors.New("from account not found")
	ErrTransferToNotFound      = errors.New("to account not found")
	ErrTransferAmountError     = errors.New("error in amount value")
	ErrTransferNoMoneyError    = errors.New("no enough money")
	ErrTransferSelfToSelfError = errors.New("disable transfer to self account")

	ErrPaymentsListNotFound         = errors.New("account not found")
	ErrPaymentsListOffsetLimitError = errors.New("error in offset, limit params")

	ErrAccountsListOffsetLimitError = errors.New("error in offset, limit params")
)

func (s Service) CreateAccount(ctx context.Context, name entity.AccountName) (entity.AccountName, error) {
	a, err := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "NewAccount()", "error", err)
		return "", ErrInService
	}
	if err = a.Validate(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Validate()", "error", err)
		return "", ErrCreateAccountInvalidName
	}
	if err = a.Register(name); err != nil {
		_ = s.logger.Log("service", "CreateAccount", "func", "Register()", "error", err)
		if a.Find(name) != err { // duplicate
			return "", ErrCreateAccountDuplicate
		}
		return "", ErrCreateAccount
	}
	return a.Name, nil
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

func (s Service) Transfer(ctx context.Context, from entity.AccountName, to entity.AccountName, amount float64) (*PaymentEntity, error) {
	aFrom, err := entity.NewAccount()
	aTo, _ := entity.NewAccount()
	if err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "NewAccount()", "error", err)
		return nil, ErrInService
	}
	if from == to {
		return nil, ErrTransferSelfToSelfError
	}
	if err = aFrom.Find(from); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return nil, ErrTransferFromNotFound
	}
	if err = aTo.Find(to); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Find()", "error", err)
		return nil, ErrTransferToNotFound
	}
	if amount <= 0 {
		return nil, ErrTransferAmountError
	}

	if aFrom.Balance < amount {
		return nil, ErrTransferNoMoneyError
	}

	if _, err = aFrom.Transfer(to, amount); err != nil {
		_ = s.logger.Log("service", "Transfer", "func", "Transfer()", "error", err)
		return nil, ErrInService
	}
	return &PaymentEntity{
		Account:   aFrom.Name,
		ToAccount: aTo.Name,
		Amount:    amount,
		Direction: PaymentDirectionOutgoing,
	}, nil
}

func (s Service) PaymentsList(ctx context.Context, name entity.AccountName, offset, limit int64) ([]PaymentEntity, error) {
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

	return convertPaymentDomainEntityToServiceEntity(lst, a)
}

func (s Service) AllPaymentsList(ctx context.Context, offset, limit int64) ([]PaymentEntity, error) {
	if offset < 0 {
		return nil, ErrPaymentsListOffsetLimitError
	}

	lst, err := entity.PaymentsList(nil, offset, limit)
	if err != nil {
		_ = s.logger.Log("service", "AllPaymentsList", "func", "List()", "error", err)
		return nil, ErrInService
	}

	return convertPaymentDomainEntityToServiceEntity(lst, nil)
}

func (s Service) AccountsList(ctx context.Context, offset, limit int64) ([]AccountEntity, error) {
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

	return convertAccountDomainEntityToServiceEntity(lst)
}

// convert response
func convertAccountDomainEntityToServiceEntity(lst []entity.Account) ([]AccountEntity, error) {
	var res []AccountEntity
	for _, a := range lst {
		res = append(res, AccountEntity{
			Id:       a.Name,
			Balance:  a.Balance,
			Currency: a.Currency,
		})
	}
	return res, nil
}

// convert response
func convertPaymentDomainEntityToServiceEntity(lst []entity.Payment, a *entity.Account) ([]PaymentEntity, error) {
	var res []PaymentEntity
	for _, p := range lst {
		// for each payment
		toAccount, err := entity.NewAccount()
		if err != nil {
			return nil, err
		}
		if err = toAccount.Get(entity.AccountID(p.ToID)); err != nil {
			toAccount = nil
		}

		fromAccount, err := entity.NewAccount()
		if err != nil {
			return nil, err
		}
		if err = fromAccount.Get(entity.AccountID(p.FromID)); err != nil {
			fromAccount = nil
		}

		to := entity.AccountName("")
		if toAccount != nil {
			to = toAccount.Name
		}
		from := entity.AccountName("")
		if fromAccount != nil {
			from = fromAccount.Name
		}
		if a != nil {
			// if defined account for witch getting payments, else all payments will be "outgoing"

			// define direction
			direction := PaymentDirectionOutgoing
			if entity.AccountID(p.FromID) != a.ID {
				from, to = to, from
				direction = PaymentDirectionIncoming
			}

			res = append(res, PaymentEntity{
				Account:   from,
				ToAccount: to,
				Amount:    p.Amount,
				Direction: direction,
			})
		} else {

			// all payments will be "outgoing"
			res = append(res, PaymentEntity{
				Account:   from,
				ToAccount: to,
				Amount:    p.Amount,
				Direction: PaymentDirectionOutgoing,
			})
		}
	}
	return res, nil
}
