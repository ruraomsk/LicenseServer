package main

import (
	"flag"
	"fmt"
	"github.com/JanFant/LicenseServer/internal/app/apiserver"
	"github.com/JanFant/LicenseServer/internal/app/db"
	"github.com/JanFant/easyLog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JanFant/LicenseServer/internal/app/config"
)

func init() {
	var configPath string
	//Начало работы, загружаем настроечный файл
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
	config.GlobalConfig = config.NewConfig()
	if _, err := toml.DecodeFile(configPath, &config.GlobalConfig); err != nil {
		fmt.Println("Can't load config file - ", err.Error())
		os.Exit(1)
	}
}

func main() {
	if err := easyLog.Init(config.GlobalConfig.LogPath); err != nil {
		fmt.Println("Error opening logger subsystem ", err.Error())
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		easyLog.Error.Println("|Message: Error open DB", err.Error())
		fmt.Println("Error open DB ", err.Error())
		return
	}
	defer dbConn.Close()

	fmt.Println("Server started...")
	easyLog.Info.Println("|Message: Server started...")
	var serverConf = apiserver.ServerConf{DB: dbConn, Port: config.GlobalConfig.ServerPort}
	apiserver.StartServer(serverConf)
}
