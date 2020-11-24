package custMain

import (
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HubTest(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		u.SendRespond(c, u.Message(http.StatusBadRequest, "bad socket connect"))
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan CustMess, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

}
