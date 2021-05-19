// Models implement interfaces for access to database layer
//
// this model can used diffrent engine of database
// for each of the engine need to define its implementation
package models

import (
	"coinswallet/pkg/wallet/models/driver"
	"fmt"
)

// interface defined account model for storage
type Account interface {
	// ID return id of wallet account
	ID() int64
	// Name return name of wallet account
	Name() string
	// Balance return balance of wallet
	Balance() float64

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
	List(offset, limit int64) ([]string, error)
}

//
// AccountFactory create model instance using dbDriver
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
