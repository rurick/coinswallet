// repository implement interfaces for access to database layer
//
// this model can used diffrent driver of database
// for each of a database need to define its implementation in driver
package repository

import (
	"fmt"
	"time"

	"coinswallet/internal/domain/wallet/repository/driver"
)

// interface defined payment repository for storage
type Payment interface {
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

	// List - return list of payments for account with accountID
	List(accountID, offset, limit int64) ([]interface{}, error)
	// ListAll - return list of all payments
	ListAll(offset, limit int64) ([]interface{}, error)
}

// create repository instance using current DBEngine
func PaymentFactory(dbDriver string) (Payment, error) {
	switch dbDriver {
	case "postgresql":
		if err := driver.PgSQLInit(); err != nil {
			return nil, err
		}
		return &driver.PgSqlPayment{}, nil
	default:
		return nil, fmt.Errorf("unknown database engine: %s", dbDriver)
	}
}
