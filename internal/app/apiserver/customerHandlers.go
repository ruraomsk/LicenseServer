package apiserver

import (
	"github.com/JanFant/LicenseServer/internal/model/customer"
	"github.com/gin-gonic/gin"
	"net/http"
)

var createCustomer = func(c *gin.Context) {
	var newCustomer customer.Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	resp := newCustomer.Create()
	c.JSON(resp.Code, resp.Obj)
}

var allCustomers = func(c *gin.Context) {
	resp := customer.GetAllCustomers()
	c.JSON(resp.Code, resp.Obj)
}
