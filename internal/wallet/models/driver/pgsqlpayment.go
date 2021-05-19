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

// PaymentsList - return list of payments for account with accountID
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// Important! When any fields will be added into table, then need to add one in to SELECT query
func PaymentsList(accountID, offset, limit int64) ([]PgSqlPayment, error) {
	rows, err := dbPool.Query(dbContext, `
		SELECT id, from, to, amount, date 
		FROM $1 
		WHERE
			from = $2 OR to = $2
		ORDER BY id
		OFFSET $3
		LIMIT $4`, paymentTableName, accountID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []PgSqlPayment
	for rows.Next() {
		pg := PgSqlPayment{}
		if err := rows.Scan(
			&pg.id, &pg.fromID, &pg.toID, &pg.amount, &pg.date); err != nil {
			return nil, err
		}
		res = append(res, pg)
	}
	return res, nil
}

// AllPaymentsList - return list of payments
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// Important! When any fields will be added into table, then need to add one in to SELECT query
func AllPaymentsList(offset, limit int64) ([]PgSqlPayment, error) {
	rows, err := dbPool.Query(dbContext, `
		SELECT id, from, to, amount, date 
		FROM $1 
		ORDER BY id
		OFFSET $3
		LIMIT $4`, paymentTableName, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []PgSqlPayment
	for rows.Next() {
		pg := PgSqlPayment{}
		if err := rows.Scan(
			&pg.id, &pg.fromID, &pg.toID, &pg.amount, &pg.date); err != nil {
			return nil, err
		}
		res = append(res, pg)
	}
	return res, nil
}
