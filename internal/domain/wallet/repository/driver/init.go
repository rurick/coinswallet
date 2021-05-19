package driver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		logger.Warning("[Wallet][getConfiguration]", err)
	}
	return configuration{
		DBName: os.Getenv("PGSQL_NAME"),
		DBUser: os.Getenv("PGSQL_USER"),
		DBHost: os.Getenv("PGSQL_HOST"),
		DBPort: os.Getenv("PGSQL_PORT"),
		DBPass: os.Getenv("PGSQL_PASS"),
	}
}

// Init - initialisation of PgSQL driver. Connect to database, checking for existing of table
// On success Set module variables dbPool,dbContext,dbCancelFunc
func PgSQLInit() (err error) {
	once.Do(func() {
		// This block will run once, when Init called first time

		config := getConfiguration()

		logger.Info("Wallet pgsql driver. Connecting to database...")
		dbContext, dbCancelFunc = context.WithCancel(context.Background())
		connStr := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
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
			}).Error("[Wallet][PgSQLInit]Unable to connect to database: ", err)
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
	if err == nil && dbPool == nil {
		err = errors.New("there's no connection to database")
	}
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
				CONSTRAINT accounts_pk PRIMARY KEY (id),
				CONSTRAINT accounts_name UNIQUE (name)
			);`
	if _, err := dbPool.Exec(dbContext, sql); err != nil {
		return fmt.Errorf("[Wallet] Can't create table accounts: %v", err)
	}

	return nil
}

// pgCreateTable is internal function used for initialisation of database
func pgCreatePaymentsTable() error {
	sql := `CREATE TABLE public.payments
			(
				id bigserial NOT NULL,
				"from" bigint NULL,
				"to" bigint NOT NULL,
				amount double precision NOT NULL DEFAULT 0,
				date timestamp with time zone NOT NULL DEFAULT now(),
				CONSTRAINT payments_pk PRIMARY KEY (id)				
			);
			CREATE INDEX payments_from_to_idx
				ON public.payments USING btree
				("from" ASC NULLS LAST, "to" ASC NULLS LAST);
			`
	if _, err := dbPool.Exec(dbContext, sql); err != nil {
		return fmt.Errorf("[Wallet] Can't create table payments: %v", err)
	}
	return nil
}
