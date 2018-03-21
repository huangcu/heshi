package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"heshi/errors"
	"net/http"
	"strings"
	"util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID                  string  `json:"id"`
	Username            string  `json:"username"`
	Cellphone           string  `json:"cellphone"`
	Email               string  `json:"email"`
	Password            string  `json:"-"`
	UserType            string  `json:"user_type"`
	RealName            string  `json:"real_name"`
	WechatID            string  `json:"wechat_id"`
	WechatName          string  `json:"wechat_name"`
	WechatQR            string  `json:"wechat_qr"`
	Address             string  `json:"address"`
	AdditionalInfo      string  `json:"additional_info"`
	RecommendedBy       string  `json:"recommended_by"`
	InvitationCode      string  `json:"invitation_code"`
	UserLevel           int     `json:"user_level"`
	UserDiscount        float64 `json:"user_discount"`
	Point               int     `json:"point"`
	TotalPurchaseAmount float64 `json:"total_purchase_amount"`
	Icon                string  `json:"icon"`
	Admin               Admin   `json:"admin"`
	Agent               Agent   `json:"agent"`
	// CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

func newAdminAgentUser(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	userType := strings.ToUpper(c.PostForm("user_type"))
	if userType != AGENT && userType != ADMIN {
		vemsgUserUsertypeNotValid.Message = fmt.Sprintf("user type can only be %s or %s", ADMIN, AGENT)
		c.JSON(http.StatusOK, vemsgUserUsertypeNotValid)
		return
	}

	nu := User{
		ID:             newV4(),
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		UserType:       userType,
		RealName:       c.PostForm("real_name"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	if userType == AGENT {
		nu.Agent = Agent{
			LevelStr:    c.PostForm("level"),
			DiscountStr: c.PostForm("discount"),
			CreatedBy:   adminID,
		}
		if vemsg, err := nu.prevalidateNewAgent(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		} else if len(vemsg) != 0 {
			c.JSON(http.StatusOK, vemsg)
			return
		}
		// q := nu.composeInsertQuery()
		// if _, err := dbExec(q); err != nil {
		// 	c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		// 	return
		// }
		// if err := a.newAgent(); err != nil {
		// 	c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		// 	return
		// }
		err := dbTransact(db, func(tx *sql.Tx) error {
			q := nu.composeInsertQuery()
			traceSQL(q)
			if _, err := tx.Exec(q); err != nil {
				return err
			}
			q = fmt.Sprintf(`INSERT INTO agents (user_id, level, discount, created_by) VALUES 
											(%s', '%d', '%d', '%s')`, nu.ID, nu.Agent.Level, nu.Agent.Discount, nu.Agent.CreatedBy)
			traceSQL(q)
			if _, err := tx.Exec(q); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, nu)
	}

	if userType == ADMIN {
		nu.Admin = Admin{
			LevelStr:   c.PostForm("level"),
			WechatKefu: c.PostForm("wechat_kefu"),
			CreatedBy:  adminID,
		}
		if vemsg, err := nu.prevalidateNewAdmin(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		} else if len(vemsg) != 0 {
			c.JSON(http.StatusOK, vemsg)
			return
		}
		// q := nu.composeInsertQuery()
		// if _, err := dbExec(q); err != nil {
		// 	c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		// 	return
		// }
		// if err := nu.newAdmin(); err != nil {
		// 	c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		// 	return
		// }
		err := dbTransact(db, func(tx *sql.Tx) error {
			q := nu.composeInsertQuery()
			traceSQL(q)
			if _, err := tx.Exec(q); err != nil {
				return err
			}
			q = fmt.Sprintf(`INSERT INTO admins (user_id, level, wechat_kefu, created_by) VALUES ('%s', '%d', '%s', '%s')`,
				nu.ID, nu.Admin.Level, nu.Admin.WechatKefu, nu.Admin.CreatedBy)
			traceSQL(q)
			if _, err := tx.Exec(q); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, nu)
	}
}

func newUser(c *gin.Context) {
	nu := User{
		ID:             newV4(),
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		UserType:       CUSTOMER,
		RealName:       c.PostForm("real_name"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	if vemsg, err := nu.validNewUser(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if len(vemsg) != 0 {
		c.JSON(http.StatusOK, vemsg)
		return
	}

	q := nu.composeInsertQuery()
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, nu.ID)
	s.Save()

	//TODO redirect to login
	// c.Redirect(http.StatusFound, redirectLogin)
	c.JSON(http.StatusOK, nu)
}

func updateAdminAgent(c *gin.Context) {
	userType := c.PostForm("user_type")
	if userType != AGENT && userType != ADMIN {
		vemsgUserUsertypeNotValid.Message = fmt.Sprintf("user type can only be %s or %s", ADMIN, AGENT)
		c.JSON(http.StatusOK, vemsgUserUsertypeNotValid)
		return
	}
	if userType == AGENT {
		updateAgent(c)
		return
	}
	if userType == ADMIN {
		updateAdmin(c)
		return
	}
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.MustGet("id").(string)
	}
	if exist, err := isUserExistByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	} else if !exist {
		c.JSON(http.StatusBadRequest, "user doesn't exist")
		return
	}
	uu := User{
		ID:        id,
		Username:  c.PostForm("username"),
		Cellphone: c.PostForm("cellphone"),
		Email:     c.PostForm("email"),
		// Password:       c.PostForm("password"), not allow update password here, call changepassword api
		UserType:       c.PostForm("user_type"),
		RealName:       c.PostForm("real_name"),
		WechatID:       c.PostForm("wechat_id"),
		WechatName:     c.PostForm("wechat_name"),
		WechatQR:       c.PostForm("wechat_qr"),
		Address:        c.PostForm("address"),
		AdditionalInfo: c.PostForm("additional_info"),
		RecommendedBy:  c.PostForm("recommended_by"),
		Icon:           c.PostForm("icon"),
	}

	//TODO validate updated user info too!!!
	//TODO what info can be updated!!
	q := uu.composeUpdateQuery()
	//TODO admin,agent update!!!!
	// var userType string
	// switch userType {
	// case "admin":
	// case "agent":
	// default:
	// }
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusBadRequest, errors.GetMessage(err))
		return
	}

	c.JSON(http.StatusOK, uu)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.MustGet("id").(string)
	}

	userType, err := getUserType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	q := selectUserQuery(id)

	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	us, err := composeUser(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if us == nil {
		vemsgUserNotExist.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, vemsgUserNotExist)
		return
	}

	if userType == CUSTOMER {
		c.JSON(http.StatusOK, us[0])
		return
	}
	if userType == ADMIN {
		a, err := getAdmin(id)
		if err != nil {
			vemsgUserNotExist.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
			c.JSON(http.StatusOK, vemsgUserNotExist)
			return
		}
		us[0].Admin = *a
		c.JSON(http.StatusOK, us[0])
		return
	}
	if userType == AGENT {
		a, err := getAgent(id)
		if err != nil {
			vemsgUserNotExist.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
			c.JSON(http.StatusOK, vemsgUserNotExist)
			return
		}
		us[0].Agent = *a
		c.JSON(http.StatusOK, us[0])
		return
	}
}

func getAllUsers(c *gin.Context) {
	q := selectUserQuery("")
	rows, err := dbQuery(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	defer rows.Close()

	us, err := composeUser(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if us == nil {
		vemsgUserNotExist.Message = "Fail to find users"
		c.JSON(http.StatusOK, vemsgUserNotExist)
		return
	}
	c.JSON(http.StatusOK, us)
}

//TODO check return row number
func disableUser(c *gin.Context) {
	uid := c.Param("id")
	q := "UPDATE users SET status='disabled' WHERE id=?"
	if _, err := dbExec(q, uid); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

func changePassword(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.MustGet("id").(string)
	}
	q := fmt.Sprintf(`SELECT password FROM users where id='%s'`, id)

	var password string
	if err := dbQueryRow(q).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, vemsgUserNotExist)
			return
		}
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	oPass := c.PostForm("old_password")
	if !util.IsPassOK(oPass, password) {
		c.JSON(http.StatusOK, vemsgLoginErrorPassword)
		return
	}
	nPass := c.PostForm("new_password")
	q = fmt.Sprintf(`update users set password='%s',updated_at=(CURRENT_TIMESTAMP) where id='%s'`, util.Encrypt(nPass), id)
	if _, err := dbExec(q); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	// TODO should relogin, clean session ??
	// s := sessions.Default(c)
	// s.Delete(USER_SESSION_KEY)
	// s.Delete(ADMIN_KEY)
	// s.Save()
	c.JSON(http.StatusOK, "PASSWORD changed!")
}

func composeUser(rows *sql.Rows) ([]User, error) {
	var id, userType, icon, invitationCode string
	var username, cellphone, email, realName, recommandedBy sql.NullString
	var wechatID, wechatName, wechatQR, address, additionalInfo sql.NullString
	var level, discount, point int
	var totalPurchaseAmount float64

	var us []User
	for rows.Next() {
		if err := rows.Scan(&id, &username, &cellphone, &email, &realName, &userType, &wechatID,
			&wechatName, &wechatQR, &address, &additionalInfo, &recommandedBy, &invitationCode,
			&level, &discount, &point, &totalPurchaseAmount, &icon); err != nil {
			return nil, err
		}
		u := User{
			ID:                  id,
			Username:            username.String,
			Cellphone:           cellphone.String,
			Email:               email.String,
			RealName:            realName.String,
			UserType:            userType,
			WechatID:            wechatID.String,
			WechatName:          wechatName.String,
			WechatQR:            wechatID.String,
			Address:             address.String,
			AdditionalInfo:      additionalInfo.String,
			RecommendedBy:       recommandedBy.String,
			InvitationCode:      invitationCode,
			UserLevel:           level,
			UserDiscount:        float64(discount) / 100,
			Point:               point,
			TotalPurchaseAmount: totalPurchaseAmount,
			Icon:                icon,
		}
		us = append(us, u)
	}
	return us, nil
}

func selectUserQuery(id string) string {
	q := `SELECT id,username,cellphone,email,real_name,user_type,wechat_id,
	wechat_name,wechat_qr,address,additional_info,recommended_by,invitation_code,
	level,discount,point,total_purchase_amount,icon FROM users`

	if id != "" {
		q = fmt.Sprintf("%s WHERE status='active' AND id='%s'", q, id)
	}
	return q
}

func getUserByID(id string) (string, error) {
	userType, err := getUserType(id)
	if err != nil {
		return "", err
	}
	q := selectUserQuery(id)

	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	us, err := composeUser(rows)
	if err != nil {
		return "", err
	}

	if userType == CUSTOMER {
		bs, err := json.Marshal(us[0])
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}
	if userType == ADMIN {
		a, err := getAdmin(id)
		if err != nil {
			return "", err
		}
		us[0].Admin = *a
		bs, err := json.Marshal(us[0])
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}
	if userType == AGENT {
		a, err := getAgent(id)
		if err != nil {
			return "", err
		}
		us[0].Agent = *a
		bs, err := json.Marshal(us[0])
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}
	return "", errors.Newf("invalid user type: %s", userType)
}
