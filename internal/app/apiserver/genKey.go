package apiserver

import (
	u "github.com/JanFant/LicenseServer/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var genKey = func(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"key": u.GenerateRandomKey(512)})
}
