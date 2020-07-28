package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code int
	Obj  map[string]interface{}
}

//Message создает map для ответа пользователю
func Message(code int, message string) Response {
	return Response{Code: code, Obj: map[string]interface{}{"message": message}}
}

//SendRespond формирует ответ пользователю записывает
func SendRespond(c *gin.Context, resp Response) {
	c.JSON(resp.Code, resp.Obj)
}
