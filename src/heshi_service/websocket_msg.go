package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type message struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"username"`
	ServeID     string    `json:"serve_id"`
	Content     string    `json:"content"`
	MsgType     string    `json:"msg_type"`
	CreatedAt   time.Time `json:"created_at"`
	MessageType int       `json:"-"`
}

func retrieveMsg(c *gin.Context) {

}

func (m *message) saveMsg() error {
	q := fmt.Sprintf(`INSERT INTO messages (user_id,serve_id, content) 
	VALUES ('%s','%s','%s')`, m.UserID, m.ServeID, m.Content)
	if m.MsgType != "" {
		q = fmt.Sprintf(`INSERT INTO messages (user_id,serve_id, msg_type, content) 
	VALUES ('%s','%s', '%s', '%s')`, m.UserID, m.ServeID, m.MsgType, m.Content)
	}
	_, err := dbExec(q)
	return err
}
