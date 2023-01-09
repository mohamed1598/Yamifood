package dbdriver

import (
	"database/sql"
	"time"
	"yamifood/pkg/helpers"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConns = 20
const maxIdleConns = 10
const maxDbLifeTime = 5 * time.Minute

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	return db, err
}

func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	helpers.ErrorCheck(err)
	db.SetMaxOpenConns(maxOpenDbConns)
	db.SetConnMaxIdleTime(maxIdleConns)
	db.SetConnMaxLifetime(maxDbLifeTime)
	dbConn.SQL = db
	return dbConn, nil
}
