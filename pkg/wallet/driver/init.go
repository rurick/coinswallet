package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	logger "github.com/sirupsen/logrus"
)

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
func PgSQLInit() (err error) {
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
		var tn string
		row := dbPool.QueryRow(dbContext, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME='payments'")
		if err = row.Scan(&tn); err != nil {
			if err = pgCreatePaymentsTable(); err != nil {
				return
			}
		}
		row = dbPool.QueryRow(dbContext, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME='accounts'")
		if err = row.Scan(&tn); err != nil {
			if err = pgCreateAccountTable(); err != nil {
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

// pgCreateTable is internal function used for initialisation of database
func pgCreateAccountTable() error {
	sql := `CREATE TABLE public.accounts
			(
				id bigserial NOT NULL,
				name character varying(32) NOT NULL,
				balance double precision NOT NULL DEFAULT 0,
				currency character varying NOT NULL,
				CONSTRAINT wallet_pk PRIMARY KEY (id),
				CONSTRAINT wallet_name UNIQUE (name)
			);`
	if _, err := dbPool.Exec(dbContext, sql); err != nil {
		return errors.New("[Wallet] Can't create table wallet")
	}

	return nil
}

// pgCreateTable is internal function used for initialisation of database
func pgCreatePaymentsTable() error {
	sql := `CREATE TABLE public.payments
			(
				id bigserial NOT NULL,
				from bigint NOT NULL,
				to bigint NOT NULL,
				amount double precision NOT NULL DEFAULT 0,
				date timestamp with time zone NOT NULL DEFAULT now(),
				CONSTRAINT wallet_pk PRIMARY KEY (id)				
			);
			CREATE INDEX from_to_idx
				ON public.payments USING btree
				(from ASC NULLS LAST, to ASC NULLS LAST);
			`
	if _, err := dbPool.Exec(dbContext, sql); err != nil {
		return errors.New("[Wallet] Can't create table payments")
	}
	return nil
}
