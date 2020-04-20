package apiserver

import (
	"github.com/JanFant/LicenseServer/internal/model/license"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var createLicense = func(c *gin.Context) {
	var newLicense license.License
	if err := c.ShouldBindJSON(&newLicense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id : cannot be blank"})
		return
	}
	resp := newLicense.CreateLicense(id)
	c.JSON(resp.Code, resp.Obj)
}

var clientInfo = func(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id : cannot be blank"})
		return
	}
	resp := license.GetAllLicenseInfo(id)
	c.JSON(resp.Code, resp.Obj)
}