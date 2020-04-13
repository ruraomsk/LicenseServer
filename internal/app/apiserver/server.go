package apiserver

import (
	"fmt"
	"github.com/JanFant/TLServer/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func StartServer(db *sqlx.DB, port string) {
	router := gin.Default()

	router.LoadHTMLGlob("./web/html/**")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inDeveloping.html", gin.H{"title": "inDevelop"})
	})

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Gin Hello")
	})

	router.StaticFS("/dir", http.Dir("./logfiles"))
	router.StaticFS("/static", http.Dir("./web/resources"))

	if err := router.Run(port); err != nil {
		logger.Error.Println("|Message: Error start server ", err.Error())
		fmt.Println("Error start server ", err.Error())
	}
}
