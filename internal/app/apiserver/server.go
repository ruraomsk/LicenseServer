package apiserver

import (
	"bufio"
	"fmt"
	"github.com/JanFant/LicenseServer/internal/app/config"
	"github.com/JanFant/LicenseServer/internal/sockets/custMain"
	"github.com/JanFant/easyLog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strings"
	"time"
)

//ServerConf конфигурация сервера
type ServerConf struct {
	DB   *sqlx.DB
	Port string
}

//StartServer запуск сервера
func StartServer(conf ServerConf) {
	hub := custMain.NewHub()
	go hub.Run()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	setLogFile()

	router := gin.Default()

	router.Use(cors.Default())

	router.LoadHTMLGlob("./web/html/**")

	router.GET("/genKey", genKey)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "custom.html", nil)
	})

	//router.GET("/custMain", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "custom.html", nil)
	//})
	router.GET("/custMain/ws", func(c *gin.Context) {
		custMain.HubTest(c, hub)
	})

	//----------------------

	router.POST("/createCustomer", createCustomer)
	router.POST("/deleteCustomer", deleteCustomer)
	router.POST("/updateCustomer", updateCustomer)

	router.GET("/client", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", nil)
	})
	router.POST("/client", clientInfo)
	router.POST("/client/createLicense", createLicense)
	router.GET("/client/createToken", createToken)

	fileServer := router.Group("/fs")
	fileServer.StaticFS("/dir", http.Dir("./logfiles"))
	fileServer.StaticFS("/res", http.Dir("./web/resources"))
	fileServer.StaticFS("/css", http.Dir("./web/css"))
	fileServer.StaticFS("/js", http.Dir("./web/js"))

	if err := router.Run(conf.Port); err != nil {
		easyLog.Error.Println("|Message: Error start server ", err.Error())
		fmt.Println("Error start server ", err.Error())
	}
}

func setLogFile() {
	path := config.GlobalConfig.GinLogPath + "/ginLog.log"
	readF, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	path2 := config.GlobalConfig.GinLogPath + "/ginLogW.log"
	writeF, _ := os.OpenFile(path2, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	scanner := bufio.NewScanner(readF)
	writer := bufio.NewWriter(writeF)
	for scanner.Scan() {
		str := scanner.Text()
		if str == "" {
			continue
		}
		splitStr := strings.Split(str, " ")
		timea, err := time.Parse("2006/01/02", splitStr[1])
		if err != nil {
			continue
		}
		if !time.Now().After(timea.Add(time.Hour * 24 * 30)) {
			_, _ = writer.WriteString(scanner.Text() + "\n")
		}
	}
	_ = writer.Flush()
	readF.Close()
	writeF.Close()

	_ = os.Remove(path)
	_ = os.Rename(path2, path)

	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	gin.DefaultWriter = file

}
