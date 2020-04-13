package db

import (
	"fmt"
	"github.com/JanFant/LicenseServer/internal/app/config"
	"github.com/JanFant/TLServer/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	licenseTable = `
	CREATE TABLE license (
		id integer PRIMARY KEY,
		numDev integer NOT NULL,
		yaKey text,
		tokenPass text,
		serverID integer,		
		endTime time
	);`
	customerTable = `
	CREATE TABLE customers (
		id integer PRIMARY KEY,
		name text,
		address text,
		servers integer[],
		phone text,
		email text
	);`
)

//ConnectDB соединение с дб
func ConnectDB() (*sqlx.DB, error) {
	db, err := sqlx.Open(config.GlobalConfig.DBConfig.Type, config.GlobalConfig.DBConfig.GetUrl())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.GlobalConfig.DBConfig.MaxOpen)
	db.SetMaxIdleConns(config.GlobalConfig.DBConfig.MaxIdle)

	_, err = db.Query(`SELECT * FROM license;`)
	if err != nil {
		fmt.Println("license table not found - created")
		logger.Info.Println("|Message: license table not found - created")
		db.MustExec(licenseTable)
	}
	_, err = db.Query(`SELECT * FROM customers;`)
	if err != nil {
		fmt.Println("customer table not found - created")
		logger.Info.Println("|Message: customer table not found - created")
		db.MustExec(customerTable)
	}
	return db, nil
}
