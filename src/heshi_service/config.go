package main

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type config struct {
	Rate      string `json:"rate"`
	CreatedBy string `json:"created_by"`
}

func getConfig(c *gin.Context) {
	var rate, createdBy string
	if err := db.QueryRow("SELECT rate,created_by FROM configs ORDER BY created_at DESC LIMIT 1").Scan(&rate, &createdBy); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, config{Rate: rate, CreatedBy: createdBy})
}

func newConfig(c *gin.Context) {
	createdBy := c.MustGet("id").(string)
	id := uuid.NewV4().String()
	q := fmt.Sprintf("INSERT INTO configs (id, rate, created_by) VALUES (%s,%s,%s)", id, createdBy, c.PostForm("rate"))
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "config set")
}
