package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"util"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

var activeConfig config

type config struct {
	Rate      float64   `json:"rate"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func getConfig(c *gin.Context) {
	var rate float64
	var createdBy string
	var createdAt time.Time
	q := "SELECT rate,created_by,created_at FROM configs ORDER BY created_at DESC LIMIT 1"
	if err := dbQueryRow(q).Scan(&rate, &createdBy, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, config{Rate: rate, CreatedBy: createdBy, CreatedAt: createdAt.Local()})
}

func newConfig(c *gin.Context) {
	// createdBy := c.MustGet("id").(string)
	createdBy := "system"
	id := uuid.NewV4().String()
	q := fmt.Sprintf("INSERT INTO configs (id, rate, created_by) VALUES ('%s','%s','%s')", id, c.PostForm("rate"), createdBy)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "config set")
}

func getAllConfigs(c *gin.Context) {
	rows, err := dbQuery("SELECT rate,created_by,created_at FROM configs ORDER BY created_at")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var cs []config
	for rows.Next() {
		var rate float64
		var createdBy string
		var createdAt time.Time
		rows.Scan(&rate, &createdBy, &createdAt)
		cs = append(cs, config{Rate: rate, CreatedBy: createdBy, CreatedAt: createdAt.Local()})
	}
	c.JSON(http.StatusOK, cs)
}

func (ac *config) getActiveConfig() {
	var rate float64
	var createdBy string
	var createdAt time.Time
	q := "SELECT rate,created_by,created_at FROM configs ORDER BY created_at DESC LIMIT 1"
	if err := dbQueryRow(q).Scan(&rate, &createdBy, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			util.Println("fail to get active config, use default config")
		}
	}
	ac.Rate = rate
	ac.CreatedBy = createdBy
	ac.CreatedAt = createdAt
}
