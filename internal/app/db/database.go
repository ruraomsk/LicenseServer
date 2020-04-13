package db

import (
	"fmt"
	"github.com/JanFant/LicenseServer/internal/app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//DBConfig конфигурации базы данных
type DBConfig struct {
	Name     string `toml:"db_name"`     //имя БД
	Password string `toml:"db_password"` //пароль доступа к БД
	User     string `toml:"db_user"`     //пользователя для обращения к бд
	Type     string `toml:"db_type"`     //тип бд
	Host     string `toml:"db_host"`     //ip сервера бд
	MaxOpen  int    `toml:"db_MaxOpen"`  //максимальное количество пустых соединений с бд
	MaxIdle  int    `toml:"db_MaxIdle"`  //максимальное количество соединенияй с бд

}

//GetUrl сформировать строку подключения к дб
func (db *DBConfig) GetUrl() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", db.Host, db.User, db.Name, db.Password)
}

var (
	licenseTable = `
	CREATE TABLE license (
		id integer PRIMARY KEY,
		numDev integer NOT NULL,
		yaKey text,
		tokenPass text,
		serverID integer,		
		endTime time
	);

	CREATE TABLE customers (
		id integer PRIMARY KEY,
		name text,
		address text,		
		phone text,
		email text,
		
	);`
)

//ConnectDB соединение с дб
func ConnectDB(dbtype, dburl string) (*sqlx.DB, error) {
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
		fmt.Println("Table here")
	} else {
		db.MustExec(licenseTable)
	}

	return db, nil
}
