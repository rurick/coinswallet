// repository implement interfaces for access to database layer

// Package repository can used different driver of database
// for each of a database need to define its implementation in driver
package repository

import (
	"fmt"

	"github.com/rurick/coinswallet/internal/domain/wallet/repository/driver"
)

// Account interface defined account repository for storage
type Account interface {
	// ID return id of wallet account
	ID() int64
	// Name return name of wallet account
	Name() string
	// Balance return balance of wallet
	Balance() float64
	// Currency return currency of wallet
	Currency() string

	// Find instance of wallet by account name
	Find(name string) error
	// Get instance of wallet by account id
	Get(id int64) error

	// Create new object in database
	Create(name string) error
	// Delete - delete wallet account
	Delete() error

	// Transfer - creating a payment form account to account with id "toID"
	Transfer(toID int64, amount float64) (int64, error)
	// Deposit - add amount to account balance
	Deposit(amount float64) (int64, error)

	// List - return list of all wallets account names
	List(offset, limit int64) ([]int64, error)
}

//
// AccountFactory create repository instance using dbDriver
func AccountFactory(dbDriver string) (Account, error) {
	switch dbDriver {
	case "postgresql":
		if err := driver.PgSQLInit(); err != nil {
			return nil, err
		}
		return &driver.PgSqlAccount{}, nil
	default:
		return nil, fmt.Errorf("unknown database engine: %s", dbDriver)
	}
}
