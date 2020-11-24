package apiserver

import (
	"github.com/JanFant/LicenseServer/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var createLicense = func(c *gin.Context) {
	var newLicense model.License
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
	resp := model.GetAllLicenseInfo(id)
	c.JSON(resp.Code, resp.Obj)
}

var createToken = func(c *gin.Context) {
	var tokenLicense model.License
	clientID, err := strconv.Atoi(c.Query("client"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id : cannot be blank"})
		return
	}
	tokenID, err := strconv.Atoi(c.Query("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id : cannot be blank"})
		return
	}
	resp := tokenLicense.CreateToken(clientID, tokenID)
	c.JSON(resp.Code, resp.Obj)
}
