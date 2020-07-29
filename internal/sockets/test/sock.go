package test

import (
	"encoding/json"
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	{
		resp := newCustomerMess(typeCustInfo, nil)
		resp.Data[typeCustInfo] = customer.GetAllCustomers()
		raw, _ := json.Marshal(resp)
		c.send <- raw
	}

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		//ну отправка и отправка
		c.hub.broadcast <- p
	}
}

func (c *Client) writePump() {
	pingTick := time.NewTicker(pingPeriod)
	defer func() {
		pingTick.Stop()
	}()
	for {
		select {
		case mess, ok := <-c.send:
			{
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "канал был закрыт"))
					return
				}

				//обычная отправка для моих нужд нужно заменить на json (наверно)
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				_, _ = w.Write(mess)

				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					_, _ = w.Write([]byte{'\n'})
					_, _ = w.Write(<-c.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			}
		case <-pingTick.C:
			{
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
