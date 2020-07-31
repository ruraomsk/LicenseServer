package test

import (
	"encoding/json"
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/JanFant/LicenseServer/internal/sockets"
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
	send chan CustMess
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	{
		resp := newCustomerMess(typeCustInfo, nil)
		resp.Data[typeCustInfo] = customer.GetAllCustomers()
		c.send <- resp
	}

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		//ну отправка и отправка
		typeSelect, err := sockets.ChoseTypeMessage(p)
		if err != nil {
			resp := newCustomerMess(typeError, nil)
			resp.Data["message"] = ErrorMessage{Error: errParseType}
			c.send <- resp
		}
		switch typeSelect {
		case typeCreateCustomer:
			{
				var cust customer.Customer
				_ = json.Unmarshal(p, &cust)
				err := cust.Create()
				if err != nil {
					resp := newCustomerMess(typeError, nil)
					resp.Data["message"] = ErrorMessage{Error: err.Error()}
					c.send <- resp
					continue
				}
				resp := newCustomerMess(typeCustUpdate, nil)
				resp.Data[typeCustInfo] = customer.GetAllCustomers()
				c.hub.broadcast <- resp
			}
		case typeDeleteCustomer:
			{
				var cust customer.Customer
				_ = json.Unmarshal(p, &cust)
				err := cust.Delete()
				if err != nil {
					resp := newCustomerMess(typeError, nil)
					resp.Data["message"] = ErrorMessage{Error: err.Error()}
					c.send <- resp
					continue
				}
				resp := newCustomerMess(typeCustUpdate, nil)
				resp.Data[typeCustInfo] = customer.GetAllCustomers()
				c.hub.broadcast <- resp
			}
		case typeUpdateCustomer:
			{
				var cust customer.Customer
				_ = json.Unmarshal(p, &cust)
				err := cust.Update()
				if err != nil {
					resp := newCustomerMess(typeError, nil)
					resp.Data["message"] = ErrorMessage{Error: err.Error()}
					c.send <- resp
					continue
				}
				resp := newCustomerMess(typeCustUpdate, nil)
				resp.Data[typeCustInfo] = customer.GetAllCustomers()
				c.hub.broadcast <- resp
			}
		default:
			{
				resp := newCustomerMess("type", nil)
				resp.Data["type"] = typeSelect
				c.send <- resp
			}
		}
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

				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				_ = json.NewEncoder(w).Encode(mess)

				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					_, _ = w.Write([]byte{'\n'})
					_ = json.NewEncoder(w).Encode(mess)
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
