package main

import (
	"encoding/json"
	"fmt"
	"heshi/errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clientList   = make(chan *client, 100)
	serveWSConns = make(map[string]*client)
	// serveMapList = make(map[*client]*client)
	sm sync.Map
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
	cc.conn.WriteMessage(1, []byte("connecting to service..."))
	if len(clientList) > 10 {
		cc.conn.WriteMessage(1, []byte("two many people waiting, we suggest you  try later"))
	} else if len(clientList) > 1 {
		cc.conn.WriteMessage(1, []byte(fmt.Sprintf("you have %d customer waiting ahead of you", len(clientList)-1)))
	} else {
		cc.conn.WriteMessage(1, []byte("connecting to serve..."))
	}

	messages := make(chan chatMessage, 100)
	serveConn := make(chan *websocket.Conn)
	go func() {
		fmt.Println("got message here")
		for {
			select {
			default:
				sm.Range(func(serve, cli interface{}) bool {
					if cli == cc {
						serveConn <- serve.(*client).conn
						return false
					}
					return true
				})
			}
		}
	}()

	go func() {
		serv := <-serveConn
		for {
			select {
			case cm := <-messages:
				m, err := json.Marshal(cm)
				if err != nil {
					cc.errChan <- err
					return
				}
				if err := serv.WriteMessage(cm.MessageType, m); err != nil {
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
	sm.Store(cc, nil)
	go distribute()
	defer func() {
		sm.Delete(cc)
	}()

	for {
		select {
		case <-cc.errChan:
			cliConn, ok := sm.Load(cc)
			if ok && cliConn != nil {
				cliConn.(*client).conn.WriteMessage(1, []byte("service down or error, please try connected again"))
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
			cliConn, ok := sm.Load(cc)
			if ok && cliConn != nil {
				m, err := json.Marshal(cm)
				if err != nil {
					cc.errChan <- errors.Mark(err)
					return
				}
				if err := cliConn.(*client).conn.WriteMessage(mt, m); err != nil {
					log.Println("write error:", errors.Mark(err))
					cc.conn.WriteMessage(mt, []byte("Customer may has ended service!"))
					sm.Store(cc, nil)
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
			go func() {
				for {
					select {
					case <-ticker.C:
						sm.Range(func(serve, cli interface{}) bool {
							if cli == nil {
								sm.Store(serve, clientConn)
								go func() {
									done <- true
								}()
								return false
							}
							return true
						})
					}
				}
			}()

			d := <-done
			fmt.Println("Done!", d)

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
	sm.Range(func(serve, _ interface{}) bool {
		if serve.(*client).status == "service" {
			count++
		}
		return true
	})
	return count
}
