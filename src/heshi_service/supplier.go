package main

import (
	"database/sql"
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
	Status    string `json:"status"`
}

func newSupplier(c *gin.Context) {
	var ns supplier
	if err := c.ShouldBind(&ns); err != nil {
		c.JSON(http.StatusOK, vemsgSupplierNotValid)
		return
	}
	if vemsg, err := ns.validUniqueKey(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}
	ns.ID = uuid.NewV4().String()
	q := ns.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, ns.ID)
}

func getAllSuppliers(c *gin.Context) {
	q := `SELECT id, name, prefix, connected, status FROM suppliers ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ss []supplier
	for rows.Next() {
		var id, name, prefix, connected, status string
		if err := rows.Scan(&id, &name, &prefix, &connected, &status); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		s := supplier{
			ID:        id,
			Name:      name,
			Prefix:    prefix,
			Connected: connected,
			Status:    status,
		}
		ss = append(ss, s)
	}
	if ss == nil {
		c.JSON(http.StatusOK, vemsgNotExist)
		return
	}
	c.JSON(http.StatusOK, ss)
}

func getSupplier(c *gin.Context) {
	q := fmt.Sprintf(`SELECT id, name, prefix, connected, status FROM suppliers WHERE id = '%s'`, c.Param("id"))
	var id, name, prefix, connected, status string
	if err := dbQueryRow(q).Scan(&id, &name, &prefix, &connected, &status); err != nil {
		if err == sql.ErrNoRows {
			vemsgNotExist.Message = fmt.Sprintf("supplier :%s not exist", c.Param("id"))
			c.JSON(http.StatusOK, vemsgNotExist)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	s := supplier{
		ID:        id,
		Name:      name,
		Prefix:    prefix,
		Connected: connected,
		Status:    status,
	}

	c.JSON(http.StatusOK, s)
}

//TODO better only allowed to change connected or not, name, prefix not allowed to change
func updateSupplier(c *gin.Context) {
	s := supplier{
		ID:        c.Param("id"),
		Name:      c.PostForm("name"),
		Prefix:    c.PostForm("prefix"),
		Connected: c.PostForm("connected"),
	}
	q := s.composeUpdateQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, s.ID)
}

//TODO check return row number?
func disableSupplier(c *gin.Context) {
	id := c.Param("id")
	q := "UPDATE suppliers SET status='disabled' WHERE id=?"
	if _, err := dbExec(q, id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func (s *supplier) composeInsertQuery() string {
	params := s.paramsKV()
	q := `INSERT INTO suppliers (id `
	va := fmt.Sprintf(`VALUES ('%s'`, s.ID)
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		switch v.(type) {
		case string:
			va = fmt.Sprintf("%s, '%s'", va, v.(string))
		case float64:
			va = fmt.Sprintf("%s, '%f'", va, v.(float64))
		case int:
			va = fmt.Sprintf("%s, '%d'", va, v.(int))
		case int64:
			va = fmt.Sprintf("%s, '%d'", va, v.(int64))
		}
	}
	return fmt.Sprintf("%s) %s)", q, va)
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

func (s *supplier) validUniqueKey() ([]errors.HSMessage, error) {
	var vemsgs []errors.HSMessage
	if exist, err := s.isSupplierExistByName(); err != nil {
		return nil, nil
	} else if exist {
		vemsgs = append(vemsgs, vemsgSupplierNameDuplicate)
	}
	if exist, err := s.isSupplierExistByPrefix(); err != nil {
		return nil, nil
	} else if exist {
		vemsgs = append(vemsgs, vemsgSupplierPrefixDuplicate)
	}

	return vemsgs, nil
}

func (s *supplier) isSupplierExistByName() (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM suppliers WHERE name='%s'", s.Name)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (s *supplier) isSupplierExistByPrefix() (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM suppliers WHERE prefix='%s'", s.Prefix)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
