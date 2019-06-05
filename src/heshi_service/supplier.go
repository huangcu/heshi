package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"
	"time"

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
	ns.ID = newV4()
	q := ns.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, ns.ID)
}

func getAllSuppliers(c *gin.Context) {
	q := `SELECT id, name, prefix, connected, status FROM suppliers ORDER BY created_at DESC`
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	var ss []supplier
	for rows.Next() {
		var id, name, prefix, connected, status string
		if err := rows.Scan(&id, &name, &prefix, &connected, &status); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
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
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
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
	uid := c.MustGet("id").(string)
	s := supplier{
		ID:        c.Param("id"),
		Name:      c.PostForm("name"),
		Prefix:    c.PostForm("prefix"),
		Connected: c.PostForm("connected"),
	}
	q := s.composeUpdateQueryTrack(uid)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, s.ID)
	// go newHistoryRecords(uid, "suppliers", s.ID, s.paramsKV())
}

//TODO check return row number?
func disableSupplier(c *gin.Context) {
	id := c.Param("id")
	q := fmt.Sprintf("UPDATE suppliers SET status='disabled' WHERE id='%s'", id)
	if _, err := dbExec(q); err != nil {
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}

	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, s.ID)
	return q
}

func (s *supplier) composeUpdateQueryTrack(updatedBy string) string {
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}
	newHistoryRecords(updatedBy, "suppliers", s.ID, params)
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, s.ID)
	return q
}

func (s *supplier) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})

	if s.Name != "" {
		params["name"] = strings.ToUpper(s.Name)
	}
	if s.Prefix != "" {
		params["prefix"] = strings.ToUpper(s.Prefix)
	}
	if s.Connected != "" {
		params["connected"] = strings.ToUpper(s.Connected)
	}
	return params
}

func (s *supplier) validUniqueKey() ([]errors.HSMessage, error) {
	var vemsgs []errors.HSMessage
	if s.Name == "" {
		vemsgNotValid.Message = "Name can not be empty"
		vemsgs = append(vemsgs, vemsgNotValid)
	} else {
		if exist, err := s.isSupplierExistByName(); err != nil {
			return nil, err
		} else if exist {
			vemsgs = append(vemsgs, vemsgSupplierNameDuplicate)
		}
	}
	if s.Prefix == "" {
		vemsgNotValid.Message = "Prefix can not be empty"
		vemsgs = append(vemsgs, vemsgNotValid)
	} else {
		if exist, err := s.isSupplierExistByPrefix(); err != nil {
			return nil, err
		} else if exist {
			vemsgs = append(vemsgs, vemsgSupplierPrefixDuplicate)
		}
	}

	return vemsgs, nil
}

func (s *supplier) isSupplierExistByName() (bool, error) {
	return isItemExistInDbByProperty("suppliers", "name", s.Name)
}

func (s *supplier) isSupplierExistByPrefix() (bool, error) {
	return isItemExistInDbByProperty("suppliers", "prefix", s.Prefix)
}
