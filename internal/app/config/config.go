package config

import "github.com/JanFant/LicenseServer/internal/app/db"

var GlobalConfig *Config

//DBConfig конфигурации базы данных
type Config struct {
	ServerPort string      `toml:"server_port"`
	LogPath    string      `toml:"log_path"`
	DBConfig   db.DBConfig `toml:"database"`
}

//NewConfig создание переменной конфиг
func NewConfig() *Config {
	return &Config{}
}
