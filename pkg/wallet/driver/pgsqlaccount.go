package driver

const (
	accountTableName = "accounts"
	defaultCurrency  = "usd"
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
		FROM $1 
		WHERE 
			"name" = $2 
		LIMIT 1`, accountTableName, name)
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
		FROM $1 
		WHERE 
			"id" = $2 
		LIMIT 1`, accountTableName, id)
	if err := row.Scan(
		&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
		return err
	}
	return nil
}

// Create - create a new account with name and load one in object
// this function not validate name
// Important! When any fields will be added into table, then need to add one in to INSERT query
func (pg *PgSqlAccount) Create(name string) error {
	res := dbPool.QueryRow(dbContext, `
		INSERT INTO $1 (name, balance, currency) VALUES(
		$2, $3, $4
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

// List - return list of all wallets accounts
// Wallets listed ordering by id
// offset and limit are using for set slice bound of list
// Important! When any fields will be added into table, then need to add one in to SELECT query
func AccountsList(offset, limit int64) ([]PgSqlAccount, error) {
	rows, err := dbPool.Query(dbContext, `
		SELECT id, name, balance, currency 
		FROM $1 
		ORDER BY id
		OFFSET $2
		LIMIT $3`, accountTableName, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []PgSqlAccount
	for rows.Next() {
		pg := PgSqlAccount{}
		if err := rows.Scan(
			&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
			return nil, err
		}
		res = append(res, pg)
	}
	return res, nil
}
