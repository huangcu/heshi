package main

import (
	"encoding/json"
	"fmt"
	"heshi/errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clientList   = make(chan *client, 100)
	serveWSConns = make(map[string]*client)
	serveMapList = make(map[*client]*client)
)

type client struct {
	id     string
	status string

	errChan chan error
	conn    *websocket.Conn
}

type chatMessage struct {
	MessageType int
	Username    string `json:"username"`
	Message     string `json:"message"`
}

var dFlag = false

func customerWSService(c *gin.Context) {
	uid := c.MustGet("id").(string)
	fmt.Println("customer id" + uid)
	uid = "guest"
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer wsConn.Close()

	cc := &client{
		conn:   wsConn,
		id:     uid,
		status: "queue",
	}

	defer func() {
		cc.status = "left"
	}()

	clientList <- cc
	fmt.Println(cc)
	cc.conn.WriteMessage(1, []byte("connecting to service..."))
	if len(clientList) > 10 {
		cc.conn.WriteMessage(1, []byte("two many people waiting, we suggest you  try later"))
	} else {
		cc.conn.WriteMessage(1, []byte(fmt.Sprintf("you have %d customer waiting ahead of you", len(clientList)-1)))
	}

	messages := make(chan chatMessage, 100)
	serveConn := make(chan *websocket.Conn)
	go func() {
		fmt.Println("got message here")
		for {
			select {
			case <-serveConn:
				return
			default:
				for serve, client := range serveMapList {
					if client == cc {
						serveConn <- serve.conn
						return
					}
				}
			}
		}
	}()

	go func() {
		s := <-serveConn
		for {
			select {
			case cm := <-messages:
				m, err := json.Marshal(cm)
				if err != nil {
					cc.errChan <- err
					return
				}
				if err := s.WriteMessage(cm.MessageType, m); err != nil {
					log.Println("write error:", errors.Mark(err))
					cc.errChan <- errors.Mark(err)
					return
				}
			}
		}
	}()

	for {
		select {
		case <-cc.errChan:
			if cc.conn != nil {
				cc.conn.WriteMessage(1, []byte("unknown error, please try connect again!"))
			}
			return
		default:
			mt, message, err := cc.conn.ReadMessage()
			if err != nil {
				fmt.Println("client connection closed")
				return
			}
			messages <- chatMessage{
				MessageType: mt,
				Username:    uid,
				Message:     string(message),
			}
		}
	}
}

func serveWSService(c *gin.Context) {
	uid := c.MustGet("id").(string)
	fmt.Println("serve id" + uid)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer wsConn.Close()

	cc := &client{
		conn:   wsConn,
		id:     uid,
		status: "service",
	}
	serveMapList[cc] = nil
	go distribute()
	defer func() {
		delete(serveMapList, cc)
	}()

	for {
		select {
		case <-cc.errChan:
			if serveMapList[cc] != nil {
				serveMapList[cc].conn.WriteMessage(1, []byte("service down or error, please try connected again"))
			}
			return
		default:
			mt, message, err := cc.conn.ReadMessage()
			if err != nil {
				cc.errChan <- errors.Mark(err)
				return
			}
			cm := chatMessage{
				MessageType: mt,
				Username:    uid,
				Message:     string(message),
			}
			if serveMapList[cc] != nil {
				m, err := json.Marshal(cm)
				if err != nil {
					cc.errChan <- errors.Mark(err)
					return
				}
				if err := serveMapList[cc].conn.WriteMessage(mt, m); err != nil {
					log.Println("write error:", errors.Mark(err))
					cc.conn.WriteMessage(mt, []byte("Customer may has ended service!"))
					serveMapList[cc].conn = nil
				}
			}
			// igonre message if no client
		}
	}
	// TODO find in DB base on uid, got last server people and latest a few messages
}

func distribute() {
	if dFlag {
		return
	}
	defer func() {
		dFlag = false
	}()
	dFlag = true

	for {
		select {
		case clientConn := <-clientList:
			done := make(chan bool)
			ticker := time.NewTicker(time.Second)

			for {
				select {
				case <-done:
					fmt.Println("Done!")
					break
				case <-ticker.C:
					for serve, client := range serveMapList {
						if client == nil {
							serveMapList[serve] = clientConn
							clientConn.conn.WriteMessage(1, []byte("what I can do for you?"))
							go func() {
								done <- true
							}()
							break
						}
					}
				}
			}
		default:
		}
	}
}

// 0, disable service connection
func isServiceAvaiable(c *gin.Context) {
	count := inServiceNumber()
	c.JSON(http.StatusOK, count)
}

func inServiceNumber() int {
	var count int
	if len(serveMapList) == 0 {
		count = 0
	} else {
		for serve := range serveMapList {
			if serve.status == "service" {
				count++
			}
		}
	}
	return count
}
