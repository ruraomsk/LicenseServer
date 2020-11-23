package db

import (
	"fmt"
	"github.com/JanFant/LicenseServer/internal/app/config"
	"github.com/JanFant/easyLog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	licenseTable = `
	CREATE TABLE license (
		id serial PRIMARY KEY,
		numdev integer NOT NULL,
		numacc integer NOT NULL,
		yakey text,
		tokenpass text,
		token text,	
		tech_email text[],
		endtime timestamp with time zone
	);`
	customerTable = `
	CREATE TABLE customers (
		id serial PRIMARY KEY,
		name text,
		address text,
		servers integer[],
		phone text,
		email text
	);`
	db *sqlx.DB
)

//ConnectDB соединение с дб
func ConnectDB() (*sqlx.DB, error) {
	conn, err := sqlx.Open(config.GlobalConfig.DBConfig.Type, config.GlobalConfig.DBConfig.GetUrl())
	if err != nil {
		return nil, err
	}
	db = conn
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.GlobalConfig.DBConfig.MaxOpen)
	db.SetMaxIdleConns(config.GlobalConfig.DBConfig.MaxIdle)

	_, err = db.Exec(`SELECT * FROM license;`)
	if err != nil {
		fmt.Println("license table not found - created")
		easyLog.Info.Println("|Message: license table not found - created")
		db.MustExec(licenseTable)
	}
	_, err = db.Exec(`SELECT * FROM customers;`)
	if err != nil {
		fmt.Println("customer table not found - created")
		easyLog.Info.Println("|Message: customer table not found - created")
		db.MustExec(customerTable)
	}
	return db, nil
}

//GetDB обращение к БД
func GetDB() *sqlx.DB {
	return db
}
