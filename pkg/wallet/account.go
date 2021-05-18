// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// This module provide handle of user wallet accounts
// For saving data in database used model defined in model.go
// Module can save data in different databases. For this action uses drivers, witch are defined in drivers directory

package wallet

import (
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
	m accountModel
}

// Register - Create new wallet account
func (a *Account) Register(name AccountName) error {
	if a.m == nil {
		return errors.New("Create new instance calling New() method first")
	}
	if err := a.m.Create(string(name)); err != nil {
		return err
	}
	return nil
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
	if a.m == nil {
		return errors.New("Create new instance calling New() method first")
	}
	err = a.m.Find(string(name))
	return
}

// Get  account by id
func (a *Account) Get(id AccountID) (err error) {
	if a.m == nil {
		return errors.New("Create new instance calling New() method first")
	}
	err = a.m.Get(int64(id))
	return
}

// load data from model to Account
func (a *Account) load() (err error) {
	if a.m == nil {
		return errors.New("Create new instance calling New() method first")
	}
	err = a.Get(a.ID)
	return
}

//
// New - create new instance of Account
func New() (*Account, error) {
	m, err := accountFactory()
	if err != nil {
		return nil, err
	}
	return &Account{
		m: m,
	}, nil
}
