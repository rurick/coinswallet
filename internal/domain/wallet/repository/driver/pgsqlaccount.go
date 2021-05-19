package driver

import "fmt"

const (
	defaultCurrency = "usd"
)

//
// Driver for work with PostgreSQL database
//
type PgSqlAccount struct {
	id       int64
	name     string
	balance  float64
	currency string
}

func (pg *PgSqlAccount) ID() int64 {
	return pg.id
}
func (pg *PgSqlAccount) Name() string {
	return pg.name
}
func (pg *PgSqlAccount) Currency() string {
	return pg.currency
}
func (pg *PgSqlAccount) Balance() float64 {
	return pg.balance
}

// Find - find wallet with name and load in object
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlAccount) Find(name string) error {
	row := dbPool.QueryRow(dbContext, `
		SELECT id, name, balance, currency 
		FROM accounts
		WHERE 
			"name" = $1 
		LIMIT 1`, name)
	if err := row.Scan(
		&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
		return err
	}
	return nil
}

// Get - get wallet by ID and load in object
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlAccount) Get(id int64) error {
	row := dbPool.QueryRow(dbContext, `
		SELECT id, name, balance, currency 
		FROM accounts
		WHERE 
			"id" = $1 
		LIMIT 1`, id)
	if err := row.Scan(
		&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
		return err
	}
	return nil
}

// Deposit - add amount to account balance
func (pg *PgSqlAccount) Deposit(amount float64) (int64, error) {
	tx, err := dbPool.Begin(dbContext)
	if err != nil {
		return 0, err
	}

	// update balances
	if _, err = tx.Exec(dbContext, `UPDATE accounts SET balance = balance + $1 WHERE id = $2 `,
		amount, pg.id); err != nil {
		if e := tx.Rollback(dbContext); e != nil {
			return 0, e
		}
		return 0, err
	}

	// create payment
	var paymentID int64
	row := tx.QueryRow(dbContext, `
		INSERT INTO payments ("from", "to", "amount", "date") VALUES(0, $1, $2, NOW()) RETURNING id`,
		pg.id, amount)
	if err = row.Scan(&paymentID); err != nil {
		if e := tx.Rollback(dbContext); e != nil {
			return 0, e
		}
		return 0, err
	}
	pg.clearPaymentsListCache()

	if err = tx.Commit(dbContext); err != nil {
		return 0, err
	}

	return paymentID, nil
}

// Transfer - creating a payment form account to account with id "toID"
// function check that recipient are exists and that the account balance is sufficient
func (pg *PgSqlAccount) Transfer(toID int64, amount float64) (int64, error) {
	tx, err := dbPool.Begin(dbContext)
	if err != nil {
		return 0, err
	}

	{
		// reread my balance from database in current transaction
		me := tx.QueryRow(dbContext, `SELECT balance FROM accounts WHERE "id" = $1 LIMIT 1`, pg.id)
		if err := me.Scan(&pg.balance); err != nil {
			return 0, err
		}
	}

	// check balance
	if pg.balance < amount {
		if err = tx.Rollback(dbContext); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("no enoth currency. balance: %f, need: %f", pg.balance, amount)
	}

	// check recipient
	to := PgSqlAccount{}
	if err = to.Get(toID); err != nil {
		return 0, fmt.Errorf("recipient not found: %v", err)
	}

	// update balances
	if _, err = tx.Exec(dbContext, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		amount, to.id); err != nil {
		if e := tx.Rollback(dbContext); e != nil {
			return 0, e
		}
		return 0, err
	}
	if _, err = tx.Exec(dbContext, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`,
		amount, pg.id); err != nil {
		if e := tx.Rollback(dbContext); e != nil {
			return 0, e
		}
		return 0, err
	}

	// create payment
	var paymentID int64
	row := tx.QueryRow(dbContext, `
		INSERT INTO payments ("from", "to", "amount", "date") VALUES($1, $2, $3, NOW()) RETURNING id`,
		pg.id, toID, amount)
	if err = row.Scan(&paymentID); err != nil {
		if e := tx.Rollback(dbContext); e != nil {
			return 0, e
		}
		return 0, err
	}
	pg.clearPaymentsListCache()

	if err = tx.Commit(dbContext); err != nil {
		return 0, err
	}

	pg.balance -= amount
	return paymentID, nil
}

// Create - create a new account with name and load one in object
// this function not validate name
// Important! When any fields will be added into table, then need to add one in to INSERT query
func (pg *PgSqlAccount) Create(name string) error {
	res := dbPool.QueryRow(dbContext, `
		INSERT INTO accounts (name, balance, currency) VALUES(
		$1, $2, $3
		)
		RETURNING id
	`, name, 0, defaultCurrency)

	var id int64
	if err := res.Scan(&id); err != nil {
		return err
	}
	pg.id = id
	pg.balance = 0
	pg.currency = defaultCurrency
	pg.name = name
	return nil
}

// Delete - delete wallet account
func (pg *PgSqlAccount) Delete() error {
	if _, err := dbPool.Exec(dbContext, `DELETE FROM accounts WHERE id = $1`, pg.id); err != nil {
		return err
	}
	return nil
}

// List - return list of all wallets account names
// Wallets listed ordering by id
// offset and limit are using for set slice bound of list
// if limit = -1, then no limit
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSqlAccount) List(offset, limit int64) ([]string, error) {
	sql := `SELECT name FROM accounts ORDER BY id OFFSET $1`
	if limit >= 0 {
		sql += fmt.Sprintf(` LIMIT %d`, limit)
	}
	rows, err := dbPool.Query(dbContext, sql, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		res []string
		s   string
	)
	for rows.Next() {
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

// clear list of account payments in cache
func (pg *PgSqlAccount) clearPaymentsListCache() {
	p := PgSqlPayment{}

	// clear cache for accounts payment list
	k := p._cacheListKey(pg.id)
	_ = cache.Delete(k)

	// clear cache for all payments list
	k = p._cacheListKey(-1)
	_ = cache.Delete(k)
}
