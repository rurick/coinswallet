package driver

import (
	"time"
)

const paymentTableName = "payments"

//
// Driver for work with PostgreSQL database
//
type PgSqlPayment struct {
	id     int64
	date   time.Time
	amount float64
	fromID int64
	toID   int64
}

func (pg *PgSqlPayment) ID() int64 {
	return pg.id
}
func (pg *PgSqlPayment) Date() time.Time {
	return pg.date
}
func (pg *PgSqlPayment) Amount() float64 {
	return pg.amount
}
func (pg *PgSqlPayment) From() int64 {
	return pg.fromID
}
func (pg *PgSqlPayment) To() int64 {
	return pg.toID
}
