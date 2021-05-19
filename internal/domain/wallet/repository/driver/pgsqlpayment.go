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

// Get - get payment by ID and load in object
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlPayment) Get(id int64) error {
	// check in cache
	cacheKey := pg._cacheKey(id)
	if v, ok := cache.Get(cacheKey); ok {
		*pg = v.(PgSqlPayment)
		return nil
	}

	row := dbPool.QueryRow(dbContext, `
		SELECT id, "from", "to", amount, date
		FROM payments
		WHERE 
			"id" = $1 
		LIMIT 1`, id)

	if err := row.Scan(&pg.id, &pg.fromID, &pg.toID, &pg.amount, &pg.date); err != nil {
		return err
	}
	cache.Set(cacheKey, *pg, 0)
	return nil
}

// List - return list of payments for account with accountID
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlPayment) List(accountID, offset, limit int64) ([]interface{}, error) {

	// try to get list from cache.
	cacheKey := pg._cacheListKey(accountID)
	if v, ok := cache.Get(cacheKey); ok {
		return v.([]interface{}), nil
	}

	sql := `
		SELECT id, "from", "to", amount, date 
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
	cache.Set(cacheKey, res, 0)
	return res, nil
}

// ListAll - return list of payments
// payments listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlPayment) ListAll(offset, limit int64) ([]interface{}, error) {

	// try to get list from cache. For list of all payments used cacheKey for accountID=-1
	cacheKey := pg._cacheListKey(-1)
	if v, ok := cache.Get(cacheKey); ok {
		return v.([]interface{}), nil
	}

	sql := `
		SELECT id, "from", "to", amount, date
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
	cache.Set(cacheKey, res, 0)
	return res, nil
}

// generate key for in memory cache
func (pg *PgSqlPayment) _cacheKey(id int64) string {
	return fmt.Sprintf("PgSqlPayment%d", id)
}

// generate key for list in memory cache
func (pg *PgSqlPayment) _cacheListKey(accountID int64) string {
	return fmt.Sprintf("List%d", accountID)
}
