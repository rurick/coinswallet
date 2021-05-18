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
)

// Account - wallet account
type Account struct {
	Name AccountName
}

// New - create new instanse of Account
func New() *Account {
	return &Account{}
}

// Register - Create new wallet account
func (a *Account) Register(name string) error {
	return nil
}

// Validate - validate account for available symbols
func (a *Account) Validate(name AccountName) error {
	re := regexp.MustCompile(`(?i)^[a-z\d]+$`)
	if !re.Match([]byte(name)) {
		return errors.New("wallet account name validate error")
	}
	return nil
}
