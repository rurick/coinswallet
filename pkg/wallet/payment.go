package wallet

import (
	"coinswallet/pkg/wallet/models"
	"time"
)

// ID identifier of payment. Integer value
type ID int64

// Payment - contain information about payment (transaction)
type Payment struct {
	ID     ID
	Date   time.Time
	Amount float64
	FromID int64
	ToID   int64

	// pointer to implementation of model
	m models.Payment
}
