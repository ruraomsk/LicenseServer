package apiserver

import (
	"github.com/JanFant/LicenseServer/internal/app/customer"
	"github.com/gin-gonic/gin"
	"net/http"
)

var createCustomer = func(c *gin.Context) {
	var newCustomer customer.Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if err := newCustomer.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": newCustomer})
}
