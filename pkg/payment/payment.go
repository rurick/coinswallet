package payment

import (
	"coinswallet/pkg/wallet"
	"time"
)

// ID identifier of payment. Integer value
type ID int64

// Payment - contain information about payment (transaction)
type Payment struct {
	ID     ID
	Date   time.Time
	Amount float64
	From   wallet.AccountName
	To     wallet.AccountName
}
