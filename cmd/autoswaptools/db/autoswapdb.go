package db

import (
	"context"
	"database/sql"
)

var (
	DefaultDBHost  = "127.0.0.1"
	DefaultDBPort  = "5432"
	DefaultDBUser  = "postgres"
	DefaultDBPass  = "a260312953"
	DefaultDBName  = "autoswap"
	DefaultTimeOut = 20

	DefaultEnableSSL = false
)

type CrossChainDB struct {
	ctx context.Context
	db  *sql.DB
}

// NewCrossChainDB
func NewCrossChainDB(dbHost string, dbPort string, dbUser string, dbPass string, DBName string) (*CrossChainDB, error) {
	//connect to the PostgreSql  daemon and return the *sql.DB
	db, err := Connect(dbHost, dbPort, dbUser, dbPass, DBName)
	if err != nil {
		return nil, err
	}

	return &CrossChainDB{
		db: db,
	}, nil

}
