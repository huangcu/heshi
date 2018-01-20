package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
	// CreatedAt      time.Time `json:"created_at"`
	// UpdatedAt      time.Time `json:"updated_at"`
}

func newAdminAgentUser(c *gin.Context) {
	adminID := c.MustGet("id").(string)
	userType := c.PostForm("user_type")
	if userType != AGENT && userType != ADMIN {
		VEMSG_USER_USERTYPE_NOT_VALID.Message = fmt.Sprintf("user type can only be %s or %s", ADMIN, AGENT)
		c.JSON(http.StatusOK, VEMSG_USER_USERTYPE_NOT_VALID)
		return
	}

	nu := User{
		ID:             uuid.NewV4().String(),
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
		a := Agent{
			User:        nu,
			LevelStr:    c.PostForm("level"),
			DiscountStr: c.PostForm("discount"),
			CreatedBy:   adminID,
		}
		if vemsg, err := a.prevalidateNewAgent(); err != nil {
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
											(%s', '%d', '%d', '%s')`, a.ID, a.Level, a.Discount, a.CreatedBy)
			traceSQL(q)
			if _, err := tx.Exec("q"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, a.ID)
	}

	if userType == ADMIN {
		a := Admin{
			User:       nu,
			LevelStr:   c.PostForm("level"),
			WechatKefu: c.PostForm("wechat_kefu"),
			CreatedBy:  adminID,
		}
		if vemsg, err := a.prevalidateNewAdmin(); err != nil {
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
		if err := a.newAdmin(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		c.JSON(http.StatusOK, a.ID)
	}
}

func newUser(c *gin.Context) {
	nu := User{
		ID:             uuid.NewV4().String(),
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

	c.JSON(http.StatusOK, nu.ID)
}

func updateAdminAgent(c *gin.Context) {
	userType := c.PostForm("user_type")
	if userType != AGENT && userType != ADMIN {
		VEMSG_USER_USERTYPE_NOT_VALID.Message = fmt.Sprintf("user type can only be %s or %s", ADMIN, AGENT)
		c.JSON(http.StatusOK, VEMSG_USER_USERTYPE_NOT_VALID)
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
	uu := User{
		ID:             id,
		Username:       c.PostForm("username"),
		Cellphone:      c.PostForm("cellphone"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
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

	c.JSON(http.StatusOK, uu.ID)
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
		VEMSG_USER_NOT_EXIST.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
		c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
		return
	}

	if userType == CUSTOMER {
		c.JSON(http.StatusOK, us[0])
		return
	}
	if userType == ADMIN {
		a, err := getAdmin(id)
		if err != nil {
			VEMSG_USER_NOT_EXIST.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
			c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
			return
		}
		a.User = us[0]
		c.JSON(http.StatusOK, a)
		return
	}
	if userType == AGENT {
		a, err := getAgent(id)
		if err != nil {
			VEMSG_USER_NOT_EXIST.Message = fmt.Sprintf("Fail to find user with id: %s", c.Param("id"))
			c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
			return
		}
		a.User = us[0]
		c.JSON(http.StatusOK, a)
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
		VEMSG_USER_NOT_EXIST.Message = "Fail to find users"
		c.JSON(http.StatusOK, VEMSG_USER_NOT_EXIST)
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
