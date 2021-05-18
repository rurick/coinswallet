// Models implement interfaces for access to database layer
//
// this model can used diffrent engine of database
// for each of the engine need to define its implementation
package wallet

import (
	"fmt"
	"time"

	"coinswallet/pkg/wallet/driver"
)

const dbEngine = "postgresql"

// interface defined account model for storage
type accountModel interface {
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

	// Transfer - creating a payment form account to account with id "toID"
	Transfer(toID int64, amount float64) error

	// List - return list of all wallets account names
	List(offset, limit int64) ([]string, error)
}

// interface defined payment model for storage
type paymentModel interface {
	// ID return id of payment
	ID() int64
	// Date return date and time of payment
	Date() time.Time
	// Amount return payments amount
	Amount() float64
	// From return payer account id
	From() int64
	// To return recipient account id
	To() int64
}

//
// create model instance using current DBEngine
func accountFactory() (accountModel, error) {
	switch dbEngine {
	case "postgresql":
		if err := driver.PgSQLInit(); err != nil {
			return nil, err
		}
		return &driver.PgSqlAccount{}, nil
	default:
		return nil, fmt.Errorf("Unknown database engine: %s", dbEngine)
	}
}

// create model instance using current DBEngine
func paymentFactory() (paymentModel, error) {
	switch dbEngine {
	case "postgresql":
		if err := driver.PgSQLInit(); err != nil {
			return nil, err
		}
		return &driver.PgSqlPayment{}, nil
	default:
		return nil, fmt.Errorf("Unknown database engine: %s", dbEngine)
	}
}
