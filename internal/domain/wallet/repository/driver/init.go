// postgresql driver for data manipulation for repository entities

// for minimisation queries count to database this driver use memory cache from package coinswallet/pkg/memcache

// configuration for connection to database getting from OS environments:
// # PostgreSQL connection
// PGSQL_HOST=127.0.0.1
// PGSQL_NAME=coins
// PGSQL_USER=coins
// PGSQL_PASS=coins
// PGSQL_PORT=5432
//
// # Memory cache settings (in minutes)
// CacheExpTime=10

// By default file with this environments locate in .env file at root of project
// this file load in getConfiguration function
// also this environments can be set in OS

package driver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	memorycache "coinswallet/pkg/memcache"
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
	// in memory cache
	cache *memorycache.Cache
	// package config
	config configuration
)

// configuration for database connection
type configuration struct {
	DBName string
	DBUser string
	DBHost string
	DBPort string
	DBPass string

	CacheExpTime time.Duration
}

// return configuration for database connection
func getConfiguration() configuration {
	if err := godotenv.Load(); err != nil {
		wd, _ := os.Getwd()
		logger.WithFields(logger.Fields{
			"[Wallet][getConfiguration]": err,
			"workdir":                    wd,
		}).Warning()
	}

	// cacheExpTime
	cacheExpTime := os.Getenv("CacheExpTime")
	if cacheExpTime == "" {
		cacheExpTime = "10"
	}
	cET, err := strconv.ParseInt(cacheExpTime, 10, 64)
	if err != nil {
		cET = 10
	}

	c := configuration{
		DBName:       os.Getenv("PGSQL_NAME"),
		DBUser:       os.Getenv("PGSQL_USER"),
		DBHost:       os.Getenv("PGSQL_HOST"),
		DBPort:       os.Getenv("PGSQL_PORT"),
		DBPass:       os.Getenv("PGSQL_PASS"),
		CacheExpTime: time.Duration(cET) * time.Minute,
	}
	logger.Debug(c)
	return c
}

// Init - initialisation of PgSQL driver. Connect to database, checking for existing of table
// On success Set module variables dbPool,dbContext,dbCancelFunc
func PgSQLInit() (err error) {
	once.Do(func() {
		// This block will run once, when Init called first time

		config = getConfiguration()
		cache = memorycache.New(config.CacheExpTime, config.CacheExpTime)

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
				balance numeric(22,4) NOT NULL DEFAULT 0,
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
				"from" bigint NOT NULL,
				"to" bigint NOT NULL,
				"amount" numeric(22,4) NOT NULL DEFAULT 0,
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
