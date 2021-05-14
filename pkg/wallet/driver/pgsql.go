package driver

import (
	"context"
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
	// Cansel function for context above
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

//return configuration for database connection
func getConfiguration() configuration {
	return configuration{}
}

// Init - initialisation of PgSQL driver. Connect to database, checking for existing of table
// On success Set module variables dbPool,dbContext,dbCancelFunc
func Init() error {
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
		dp, err := pgxpool.Connect(dbContext, connStr)

		dbPool = dp
		if err != nil {
			logger.WithFields(logger.Fields{
				"DBUser": config.DBUser,
				"DBPass": "********",
				"DBName": config.DBName,
				"DBHost": config.DBHost,
				"DBPort": config.DBPort,
			}).Panic("Unable to connect to database: %v", err)
		}

		// checking for existing of table
		row := dbPool.QueryRow(dbContext, "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME='%s'", tableName)
		var tn string
		if err := row.Scan(&tn); err != nil {
			logger.Panic("There is no table %v in database", tableName)
		}

		// Close db connection when context ware completed
		go func(ctx context.Context) {
			<-ctx.Done()
			dbPool.Close()
			dbCancelFunc()
			logger.Info("Wallet pgsql driver. DB connection closed")
		}(dbContext)
	})
	return nil
}

// Driver for work with PostgreSQL database
type PgSql struct {
	name string
}

func (pg *PgSql) Name() string {
	return pg.name
}
func (pg *PgSql) Find(name string) error {
	row := dbPool.QueryRow(dbContext, `
		SELECT 
			"name"
		FROM $1 
		WHERE 
			"name" = $2 
		LIMIT 1`, tableName, name)
	if err := row.Scan(&pg.name); err != nil {
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
