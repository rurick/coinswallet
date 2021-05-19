package entity

import (
	"time"

	"coinswallet/internal/domain/wallet/repository"
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
	rep repository.Payment
}

func (a *Payment) load() {
	a.ID = ID(a.rep.ID())
	a.Amount = a.rep.Amount()
	a.Date = a.rep.Date()
	a.FromID = a.rep.From()
	a.ToID = a.rep.To()
}

// Get  account by id
func (a *Payment) Get(id ID) (err error) {
	err = a.rep.Get(int64(id))
	if err == nil {
		a.load()
	}
	return
}

//
// NewPayment - create new instance of Payment
func NewPayment() (*Payment, error) {
	const dbDriver = "postgresql"

	rep, err := repository.PaymentFactory(dbDriver)
	if err != nil {
		return nil, err
	}
	return &Payment{
		rep: rep,
	}, nil
}

// PaymentsList - return list of all payments.
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// if account is nil returning list of all accounts
func PaymentsList(account *Account, offset, limit int64) ([]Payment, error) {
	p, err := NewPayment()
	if err != nil {
		return nil, err
	}
	var lst []interface{}
	if account == nil {
		if lst, err = p.rep.ListAll(offset, limit); err != nil {
			return nil, err
		}
	} else {
		if lst, err = p.rep.List(int64(account.ID), offset, limit); err != nil {
			return nil, err
		}
	}

	var res []Payment
	for _, n := range lst {
		rp := n.(repository.Payment)
		p := Payment{rep: rp} // it is able don't call NewPayment() because it was called early at begin of function so  driver was initialised
		p.load()
		res = append(res, p)
	}
	return res, nil
}
