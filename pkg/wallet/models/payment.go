// Models implement interfaces for access to database layer
//
// this model can used diffrent engine of database
// for each of the engine need to define its implementation
package models

import (
	"fmt"
	"time"

	"coinswallet/pkg/wallet/models/driver"
)

// interface defined payment model for storage
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
}

// create model instance using current DBEngine
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
