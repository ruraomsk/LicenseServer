package config

import (
	"fmt"
)

var GlobalConfig *Config

//DBConfig конфигурации базы данных
type Config struct {
	ServerPort string   `toml:"server_port"`
	LogPath    string   `toml:"log_path"`
	DBConfig   DBConfig `toml:"database"`
}

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

//NewConfig создание переменной конфиг
func NewConfig() *Config {
	return &Config{}
}
