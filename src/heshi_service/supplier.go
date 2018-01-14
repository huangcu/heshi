package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

// id VARCHAR(225) PRIMARY KEY NOT NULL,
// name VARCHAR(255) NOT NULL,
// prefix VARCHAR(8) NOT NULL,
// connected VARCHAR(5) NOT NULL,
type supplier struct {
	ID        string `json:"id"`
	Name      string `form:"name" json:"name" binding:"required"`
	Prefix    string `form:"prefix" json:"prefix" binding:"required"`
	Connected string `form:"connected" json:"connected"`
}

func newSupplier(c *gin.Context) {
	var ns supplier
	if err := c.ShouldBind(&ns); err != nil {
		c.JSON(http.StatusOK, VEMSG_SUPPLIER_NOT_VALID)
		return
	}
	ns.ID = uuid.NewV4().String()
	q := fmt.Sprintf(`INSERT INTO suppliers (id, name, prefix, connected) VALUES ('%s', '%s', '%s', '%s')`,
		ns.ID, ns.Name, ns.Prefix, ns.Connected)
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, ns.ID)
}

func getAllSuppliers(c *gin.Context) {
	q := `SELECT id, name, prefix, connected FROM suppliers ORDER BY created_at DESC`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ss []supplier
	for rows.Next() {
		var id, name, prefix, connected string
		if err := rows.Scan(&id, &name, &prefix, &connected); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		s := supplier{
			ID:        id,
			Name:      name,
			Prefix:    prefix,
			Connected: connected,
		}
		ss = append(ss, s)
	}
	if ss == nil {
		c.JSON(http.StatusOK, VEMSG_NOT_EXIST)
		return
	}
	c.JSON(http.StatusOK, ss)
}

func updateSupplier(c *gin.Context) {
	s := supplier{
		ID:        c.Param("id"),
		Name:      c.PostForm("name"),
		Prefix:    c.PostForm("prefix"),
		Connected: c.PostForm("connected"),
	}
	q := s.composeUpdateQuery()
	if _, err := db.Exec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, s.ID)
}

//TODO check return row number?
func removeSupplier(c *gin.Context) {
	id := c.Param("id")
	q := "DELETE FROM suppliers WHERE id=?"
	if _, err := db.Exec(q, id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func (s *supplier) composeUpdateQuery() string {
	params := s.paramsKV()
	q := `UPDATE suppliers SET`
	for k, v := range params {
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		}
	}

	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), s.ID)
	return q
}

func (s *supplier) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if s.Name != "" {
		params["name"] = s.Name
	}
	if s.Prefix != "" {
		params["prefix"] = s.Prefix
	}
	if s.Connected != "" {
		params["connected"] = s.Connected
	}
	return params
}
