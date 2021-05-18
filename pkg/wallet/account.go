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
	id      AccountID
	name    AccountName
	balance float64
}

// New - create new instance of Account
func New(id AccountID) *Account {
	return &Account{}
}

// Register - Create new wallet account
func (a *Account) Register(name string) error {
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
func (a *Account) Find(name AccountName) error {
	return nil
}

func (a *Account) Balance() (float64, error) {
	if a == nil {
		return 0, errors.New("Create new instance calling New() method first")
	}
	return 0, nil
}
