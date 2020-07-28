package customer

import (
	"github.com/JanFant/TLServer/logger"
	"github.com/gorilla/websocket"
)

var (
	typeError      = "error"
	typeClose      = "close"
	typeCustInfo   = "custInfo"
	typeCustUpdate = "custUpdate"
)

type CustMess struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	conn *websocket.Conn        `json:"-"`
}

func newCustomerMess(mType string, conn *websocket.Conn, data map[string]interface{}) CustMess {
	var resp CustMess
	resp.Type = mType
	resp.conn = conn
	if data != nil {
		resp.Data = data
	} else {
		resp.Data = make(map[string]interface{})
	}
	return resp
}

func (c *CustMess) send() {
	if c.Type == typeError {
		go func() {
			logger.Warning.Printf("|IP: %s |Login: %s |Resource: %s |Message: %v", c.conn.RemoteAddr(), "map socket", "/map", c.Data["message"])
		}()
	}
	writeMessage <- *c
}

//ErrorMessage структура ошибки
type ErrorMessage struct {
	Error string `json:"error"`
}
