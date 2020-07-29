package apiserver

import (
	"fmt"
	"github.com/JanFant/LicenseServer/internal/sockets/customer"
	"github.com/JanFant/LicenseServer/internal/sockets/test"
	"github.com/JanFant/TLServer/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

//ServerConf конфигурация сервера
type ServerConf struct {
	DB   *sqlx.DB
	Port string
}

//StartServer запуск сервера
func StartServer(conf ServerConf) {
	go customer.CustBroadcast()
	hub := test.NewHub()
	go hub.Run()

	router := gin.Default()

	router.Use(cors.Default())

	router.LoadHTMLGlob("./web/html/**")

	router.GET("/genKey", genKey)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "custom.html", nil)
	})
	//router.POST("/", allCustomers)

	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "custom.html", nil)
	})
	router.GET("/test/ws", func(c *gin.Context) {
		test.HubTest(c, hub)
	})

	//----------------------
	router.GET("/ws", customerEngine)
	//-------------

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
		logger.Error.Println("|Message: Error start server ", err.Error())
		fmt.Println("Error start server ", err.Error())
	}
}
