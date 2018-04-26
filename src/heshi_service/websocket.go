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
	clientList   = make(chan *socketConn, 100)
	serveWSConns = make(map[string]*socketConn)
	// serveMapList = make(map[*socketConn]*socketConn)
	sm sync.Map
)

// status: queue, in service(for customer, it is served by a serve, for serve, it is in service), left
type socketConn struct {
	id     string
	status string

	errChan chan error
	conn    *websocket.Conn
}

// type chatMessage struct {
// 	MessageType int
// 	Username    string `json:"username"`
// 	Message     string `json:"message"`
// }

var dFlag = false

func customerWSService(c *gin.Context) {
	uid := c.MustGet("id").(string)
	username, err := getUserNameByID(uid)
	if err != nil {
		log.Println(err)
		return
	}
	if username == "" {
		username = uid
	}
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

	cc := &socketConn{
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

	messages := make(chan message, 100)
	serveConn := make(chan *socketConn)
	exitChan1 := make(chan bool)
	exitChan2 := make(chan bool)

	go func() {
		for {
			select {
			case <-exitChan1:
				return
			default:
				sm.Range(func(serve, cli interface{}) bool {
					if cli == cc && cc.status == "queue" {
						cc.status = "in service"
						serveConn <- serve.(*socketConn)
						return false
					}
					return true
				})
			}
		}
	}()

	go func() {
		for {
			select {
			case <-exitChan2:
				return
			case serv := <-serveConn:
				go func() {
					for {
						select {
						case <-exitChan2:
							sm.Store(serv, nil)
							return
						case cm := <-messages:
							m, err := json.Marshal(cm)
							if err != nil {
								cc.errChan <- err
								return
							}
							if err := serv.conn.WriteMessage(cm.MessageType, m); err != nil {
								log.Println("write error:", errors.Mark(err))
								cc.errChan <- errors.Mark(err)
								return
							}
						}
					}
				}()
				return
			}
		}
	}()

	for {
		select {
		case <-cc.errChan:
			if cc.conn != nil {
				cc.conn.WriteMessage(1, []byte("unknown error, please try connect again!"))
			}
			go func() { exitChan1 <- true }()
			go func() { exitChan2 <- true }()
			return
		default:
			mt, msg, err := cc.conn.ReadMessage()
			if err != nil {
				fmt.Println("client connection closed")
				go func() { exitChan1 <- true }()
				go func() { exitChan2 <- true }()
				return
			}
			m := message{
				MessageType: mt,
				UserID:      uid,
				UserName:    username,
				Content:     string(msg),
			}
			go m.saveMsg()
			messages <- m
		}
	}
}

func serveWSService(c *gin.Context) {
	uid := c.MustGet("id").(string)
	username, err := getUserNameByID(uid)
	if err != nil {
		log.Println(err)
		return
	}
	if username == "" {
		username = uid
	}
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

	cc := &socketConn{
		conn:   wsConn,
		id:     uid,
		status: "in service",
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
				cliConn.(*socketConn).conn.WriteMessage(1, []byte("service down or error, please try connected again"))
				cliConn.(*socketConn).conn.Close()
			}
			return
		default:
			mt, msg, err := cc.conn.ReadMessage()
			if err != nil {
				fmt.Println("serve connection closed")
				cliConn, ok := sm.Load(cc)
				if ok && cliConn != nil {
					cliConn.(*socketConn).conn.WriteMessage(1, []byte("service down or error, please try connected again"))
					cliConn.(*socketConn).conn.Close()
				}
				return
			}
			cm := message{
				MessageType: mt,
				ServeID:     uid,
				Content:     string(msg),
			}
			cliConn, ok := sm.Load(cc)
			if ok && cliConn != nil {
				cm.UserID = cliConn.(*socketConn).id
				if cliConn.(*socketConn).status != "left" {
					go cm.saveMsg()
					if err := cliConn.(*socketConn).conn.WriteMessage(mt, msg); err != nil {
						log.Println("write error:", errors.Mark(err))
						cc.conn.WriteMessage(mt, []byte("Customer may has ended service!"))
						sm.Store(cc, nil)
					}
				} else {
					cc.conn.WriteMessage(mt, []byte("Customer has ended service!"))
				}
			}
			// igonre message if no client
		}
	}
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
							} else if cli.(*socketConn).status == "left" {
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
		if serve.(*socketConn).status == "in service" {
			count++
		}
		return true
	})
	return count
}
