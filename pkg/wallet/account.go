// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// This module provide handle of user wallet accounts
// For saving data in database used model defined in model.go
// Module can save data in different databases. For this action uses drivers, witch are defined in drivers directory

package wallet

import (
	"coinswallet/pkg/wallet/models"
	"errors"
	"regexp"
)

type (
	// AccountName name of user account. Available symbol are latin letters or numbers
	AccountName string
	AccountID   int64
)

// Account - wallet account
type Account struct {
	ID       AccountID
	Name     AccountName
	Balance  float64
	Currency string

	// pointer to implementation of model
	m models.Account
	// database driver
	dbDriver string
}

// Register - Create a new wallet account with zero balance
func (a *Account) Register(name AccountName) (err error) {
	err = a.m.Create(string(name))
	return
}

// Delete - delete wallet account
func (a *Account) Delete() (err error) {
	err = a.m.Delete()
	return
}

// Validate - validate account for available symbols and length (4-32)
func (a *Account) Validate(name AccountName) error {
	re := regexp.MustCompile(`(?i)^[a-z\d]{4,32}$`)
	if !re.Match([]byte(name)) {
		return errors.New("wallet account name validate error")
	}
	return nil
}

// Find find account by name
func (a *Account) Find(name AccountName) (err error) {
	err = a.m.Find(string(name))
	return
}

// Get  account by id
func (a *Account) Get(id AccountID) (err error) {
	err = a.m.Get(int64(id))
	return
}

// Transfer creating a payment form account "a" to account with id "toID"
// returning id of payment
func (a *Account) Transfer(toName AccountName, amount float64) (id int64, err error) {
	var to *Account

	if to, err = New(); err != nil {
		return
	}
	if err = to.Find(toName); err != nil {
		return
	}

	id, err = a.m.Transfer(int64(to.ID), amount)
	return
}

// Deposit - add amount to account balance.
// returning id of payment
func (a *Account) Deposit(amount float64) (id int64, err error) {
	id, err = a.m.Deposit(amount)
	return
}

// List - return list of all wallets account names
// Wallets listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
func (a *Account) List(offset, limit int64) ([]AccountName, error) {
	lst, err := a.m.List(offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert result type
	var res []AccountName
	for _, n := range lst {
		res = append(res, AccountName(n))
	}
	return res, nil
}

//
// load data from model to Account
func (a *Account) load() (err error) {
	if a.m == nil {
		return errors.New("create new instance calling New() method first")
	}
	err = a.Get(a.ID)
	return
}

//
// New - create new instance of Account
func New() (*Account, error) {
	const dbDriver = "postgresql"

	m, err := models.AccountFactory(dbDriver)
	if err != nil {
		return nil, err
	}
	return &Account{
		m: m,
	}, nil
}
