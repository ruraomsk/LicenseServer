package apiserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var count = 0

//allCustomers обработчик запроса получения всех клиентов
var allCustomersWS = func(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer conn.Close()
	count++
	fmt.Println("Client connected: ", count)

	reader(conn, count)
}

func reader(conn *websocket.Conn, count1 int) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p), "   ", count1)

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}

	}
}
