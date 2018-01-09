package main

import (
	"net/http"
	"util"

	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) {
	t, err := util.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": t})
}

func VerifyToken(c *gin.Context) {
	t := c.PostForm("token")
	if t == "" {
		c.JSON(http.StatusForbidden, "No token")
		return
	}
	if util.VerfiyToken(t) {
		c.JSON(http.StatusOK, "Hi it's you!")
		return
	}
	c.JSON(http.StatusForbidden, "Hi you don't have the right token")
}
