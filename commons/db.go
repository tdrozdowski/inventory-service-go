package commons

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"sync"
	"time"
)

var (
	// DB is a global variable that holds the connection to the database
	sqlxDb *sqlx.DB
	once   sync.Once
)

func getDbUri() string {
	dbConnectionString := os.Getenv("DATABASE_URL")
	return dbConnectionString
}

func GetDB() *sqlx.DB {
	once.Do(func() {
		uri := getDbUri()
		db, err := pgxCreateDB(uri)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(4)
		db.SetConnMaxLifetime(time.Duration(30) * time.Minute)
		sqlxDb = db
	})

	return sqlxDb
}

func pgxCreateDB(uri string) (*sqlx.DB, error) {
	pool, err := pgxpool.New(context.Background(), uri)

	if err != nil {
		return nil, err
	}

	afterConnect := stdlib.OptionAfterConnect(func(ctx context.Context, conn *pgx.Conn) error {
		// for anything we decide we need to do after the connection is established
		return nil
	})

	pgxdb := stdlib.OpenDBFromPool(pool, afterConnect)
	return sqlx.NewDb(pgxdb, "pgx"), nil
}
