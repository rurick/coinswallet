package driver

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	logger "github.com/sirupsen/logrus"
	"sync"
)

const tableName = "accounts"

var (
	// pool of database resources
	dbPool *pgxpool.Pool
	// once used for initialisation of db connection
	once sync.Once
	// Context of all database operations
	dbContext context.Context
	// Cancel function for context above
	dbCancelFunc context.CancelFunc
)

// configuration for database connection
type configuration struct {
	DBName string
	DBUser string
	DBHost string
	DBPort string
	DBPass string
}

// return configuration for database connection
func getConfiguration() configuration {
	return configuration{}
}

// Init - initialisation of PgSQL driver. Connect to database, checking for existing of table
// On success Set module variables dbPool,dbContext,dbCancelFunc
func Init() (err error) {
	config := getConfiguration()
	once.Do(func() {
		// This block will run once, when Init called first time

		logger.Info("Wallet pgsql driver. Connecting to database...")
		dbContext, dbCancelFunc = context.WithCancel(context.Background())
		connStr := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable ",
			config.DBUser,
			config.DBPass,
			config.DBName,
			config.DBHost,
			config.DBPort,
		)
		dbPool, err = pgxpool.Connect(dbContext, connStr)
		if err != nil {
			logger.WithFields(logger.Fields{
				"DBUser": config.DBUser,
				"DBPass": "********",
				"DBName": config.DBName,
				"DBHost": config.DBHost,
				"DBPort": config.DBPort,
			}).Error("[Wallet][Init]Unable to connect to database: %v", err)
			return
		}

		// checking for existing of table
		row := dbPool.QueryRow(dbContext, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME='%s'", tableName)
		var tn string
		if err = row.Scan(&tn); err != nil {
			if err = createTable(); err != nil {
				return
			}
		}

		// Close db connection when context ware completed
		go func(ctx context.Context) {
			<-ctx.Done()
			dbPool.Close()
			dbCancelFunc()
			logger.Info("Wallet pgsql driver. DB connection closed")
		}(dbContext)
	})
	return
}

// createTable is internal function used for initialisation of database
func createTable() error {
	sql := `CREATE TABLE public.wallet
			(
				name character varying(32) NOT NULL,
				balance double precision NOT NULL DEFAULT 0,
				currency character varying NOT NULL,
				id bigserial NOT NULL,
				CONSTRAINT wallet_pk PRIMARY KEY (id),
				CONSTRAINT wallet_name UNIQUE (name)
			);`
	if _, err := dbPool.Exec(dbContext, sql); err != nil {
		return errors.New("[Wallet] Can't create table wallet")
	}
	return nil
}

//
// Driver for work with PostgreSQL database
type PgSql struct {
	id       int64
	name     string
	balance  float64
	currency string
}

func (pg *PgSql) ID() int64 {
	return pg.id
}
func (pg *PgSql) Name() string {
	return pg.name
}
func (pg *PgSql) Currency() string {
	return pg.currency
}
func (pg *PgSql) Balance() float64 {
	return pg.balance
}

// Find - find wallet with name and load in object
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSql) Find(name string) error {
	row := dbPool.QueryRow(dbContext, `
		SELECT id, name, balance, currency 
		FROM $1 
		WHERE 
			"name" = $2 
		LIMIT 1`, tableName, name)
	if err := row.Scan(
		&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
		return err
	}
	return nil
}

// Get - get wallet by ID and load in object
// Important! When any fields will be added into table, then need to add one in to SELECT query
func (pg *PgSql) Get(id int64) error {
	row := dbPool.QueryRow(dbContext, `
		SELECT id, name, balance, currency 
		FROM $1 
		WHERE 
			"id" = $2 
		LIMIT 1`, tableName, id)
	if err := row.Scan(
		&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
		return err
	}
	return nil
}

func (pg *PgSql) Save() error {
	return nil
}

func (pg *PgSql) Create() error {
	return nil
}

// List - return list of all wallets accounts
// Wallets listed ordering by id
// offset and limit are using for set slice bound of list
// Important! When any fields will be added into table, then need to add one in to SELECT query
func List(offset, limit int64) ([]PgSql, error) {
	rows, err := dbPool.Query(dbContext, `
		SELECT id, name, balance, currency 
		FROM $1 
		ORDER BY id
		OFFSET $2
		LIMIT $3`, tableName, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []PgSql
	for rows.Next() {
		pg := PgSql{}
		if err := rows.Scan(
			&pg.id, &pg.name, &pg.balance, &pg.currency); err != nil {
			return nil, err
		}
		res = append(res, pg)
	}
	return res, nil
}
