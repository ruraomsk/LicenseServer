package apiserver

import (
	"fmt"
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

	router := gin.Default()

	router.Use(cors.Default())

	router.LoadHTMLGlob("./web/html/**")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", gin.H{"title": "inDevelop"})
	})

	router.GET("/genKey", genKey)

	mainRouter := router.Group("/main")

	mainRouter.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", gin.H{"title": "inDevelop"})
	})
	mainRouter.POST("/", allCustomers)

	mainRouter.POST("/createCustomer", createCustomer)
	mainRouter.POST("/deleteCustomer", deleteCustomer)
	mainRouter.POST("/updateCustomer", updateCustomer)

	mainRouter.GET("/client", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", gin.H{"title": "inDevelop"})
	})
	mainRouter.POST("/client", clientInfo)
	mainRouter.POST("/client/createLicense", createLicense)
	mainRouter.GET("/client/createToken", createToken)

	fileServer := router.Group("/fs")
	fileServer.StaticFS("/dir", http.Dir("./logfiles"))
	fileServer.StaticFS("/static", http.Dir("./web/resources"))

	if err := router.Run(conf.Port); err != nil {
		logger.Error.Println("|Message: Error start server ", err.Error())
		fmt.Println("Error start server ", err.Error())
	}
}
