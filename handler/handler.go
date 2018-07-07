package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/livestockz/api/utils"
)

func HealthHandler(c *gin.Context) {
	utils.Ok(c, nil)
}
