package driver

import (
	"fmt"
	"time"
)

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

// List - return list of payments for account with accountID
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlPayment) List(accountID, offset, limit int64) ([]interface{}, error) {
	sql := `
		SELECT id, from, to, amount, date 
		FROM payments
		WHERE
			from = $1 OR to = $1
		ORDER BY id
		OFFSET $2`
	if limit >= 0 {
		sql += fmt.Sprintf(` LIMIT %d`, limit)
	}

	rows, err := dbPool.Query(dbContext, sql, accountID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []interface{}
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

// ListAll - return list of payments
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlPayment) ListAll(offset, limit int64) ([]interface{}, error) {
	sql := `
		SELECT id, from, to, amount, date 
		FROM payments
		ORDER BY id
		OFFSET $1`
	if limit >= 0 {
		sql += fmt.Sprintf(` LIMIT %d`, limit)
	}

	rows, err := dbPool.Query(dbContext, sql, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []interface{}
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
