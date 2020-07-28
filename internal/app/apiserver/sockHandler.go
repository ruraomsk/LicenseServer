package apiserver

import (
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/JanFant/LicenseServer/internal/sockets/customer"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var customerEngine = func(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		u.SendRespond(c, u.Message(http.StatusBadRequest, "bad socket connect"))
		return
	}
	defer conn.Close()
	customer.CustReader(conn)
}
