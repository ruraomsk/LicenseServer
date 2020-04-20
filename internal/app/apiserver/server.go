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

	mainRouter := router.Group("/main")

	mainRouter.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", gin.H{"title": "inDevelop"})
	})
	mainRouter.POST("/", allCustomers)
	mainRouter.POST("/createCustomer", createCustomer)
	mainRouter.POST("/deleteCustomer", deleteCustomer)
	mainRouter.POST("/updateCustomer", updateCustomer)

	fileServer := router.Group("/fileServer")
	fileServer.StaticFS("/dir", http.Dir("./logfiles"))
	fileServer.StaticFS("/static", http.Dir("./web/resources"))

	if err := router.Run(conf.Port); err != nil {
		logger.Error.Println("|Message: Error start server ", err.Error())
		fmt.Println("Error start server ", err.Error())
	}
}
