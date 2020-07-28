package customer

import (
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/JanFant/LicenseServer/internal/sockets"
	"github.com/gorilla/websocket"
	"reflect"
	"time"
)

var writeMessage chan CustMess
var poolConnect map[*websocket.Conn]bool

func CustReader(conn *websocket.Conn) {
	poolConnect[conn] = true
	{
		resp := newCustomerMess(typeCustInfo, conn, nil)
		resp.Data[typeCustInfo] = customer.GetAllCustomers()
		resp.send()
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			resp := newCustomerMess(typeClose, conn, nil)
			resp.send()
			return
		}
		typeSelect, err := sockets.ChoseTypeMessage(p)
		if err != nil {
			resp := newCustomerMess(typeError, conn, nil)
			resp.Data["message"] = "Нет такого типа в запросах"
			resp.send()
		}
		switch typeSelect {
		case "nothing":
			{

			}
		}
	}
}

func CustBroadcast() {
	writeMessage = make(chan CustMess, 5)
	poolConnect = make(map[*websocket.Conn]bool)
	ticker := time.NewTicker(time.Second * 10)
	oldCust := customer.GetAllCustomers()
	for {
		select {
		case <-ticker.C:
			{
				if len(poolConnect) > 0 {
					var temp []customer.Customer
					newCust := customer.GetAllCustomers()
					for _, nCust := range newCust {
						flagN := true
						for _, oCust := range oldCust {
							if nCust.Name == oCust.Name {
								flagN = false
								if !reflect.DeepEqual(nCust, oCust) {
									temp = append(temp, nCust)
								}
							}
						}
						if flagN {
							temp = append(temp, nCust)
						}
					}
					oldCust = newCust
					if len(temp) > 0 {
						resp := newCustomerMess(typeCustUpdate, nil, nil)
						resp.Data[typeCustUpdate] = temp
						for conn := range poolConnect {
							_ = conn.WriteJSON(resp)
						}
					}
				}
			}
		case msg := <-writeMessage:
			{
				switch msg.Type {
				case typeClose:
					{
						delete(poolConnect, msg.conn)
					}
				default:
					{
						_ = msg.conn.WriteJSON(msg)
					}
				}
			}
		}
	}
}
