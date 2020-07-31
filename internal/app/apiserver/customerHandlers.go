package apiserver

import (
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/gin-gonic/gin"
	"net/http"
)

////allCustomers обработчик запроса получения всех клиентов
//var allCustomers = func(c *gin.Context) {
//	resp := customer.GetAllCustomers()
//	c.JSON(resp.Code, resp.Obj)
//}

//createCustomer обработчик создания клиента
var createCustomer = func(c *gin.Context) {
	var newCustomer customer.Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	//resp := newCustomer.Create()
	//c.JSON(resp.Code, resp.Obj)
}

//deleteCustomer обработчик удаления клиента
var deleteCustomer = func(c *gin.Context) {
	var delCustomer customer.Customer
	if err := c.ShouldBindJSON(&delCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	//resp := delCustomer.Delete()
	//c.JSON(resp.Code, resp.Obj)
}

//updateCustomer обработчик обновления данных клиента
var updateCustomer = func(c *gin.Context) {
	var updateCustomer customer.Customer
	if err := c.ShouldBindJSON(&updateCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	//resp := updateCustomer.Update()
	//c.JSON(resp.Code, resp.Obj)
}
