package main

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

//微信公众号平台Callback - 接入验证
func checkSignature(c *gin.Context) {
	token := "BEYOU_SIHUI"
	echostr := c.Query("echostr")
	//接入验证
	if echostr != "" {
		signature := c.Query("signature")
		tmpArr := []string{token, c.Query("timestamp"), c.Query("nonce")}
		sort.Strings(tmpArr)
		h := sha1.New()
		_, err := h.Write([]byte(strings.Join(tmpArr, "")))
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		bs := h.Sum(nil)
		s := fmt.Sprintf("%x", bs)

		if s == signature {
			//accessed
			c.String(http.StatusOK, echostr)
			return
		}
	}

	body, err := c.Request.GetBody()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	bs, err := ioutil.ReadAll(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var msg MixedMsg
	err = xml.Unmarshal(bs, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if msg.MsgType == "event" {
		switch msg.EventType {
		case EventTypeSubscribe:
		case EventTypeUnsubscribe:
		case EventTypeScan:
			if err := redisClient.Set(msg.EventKey, msg.FromUserName, 0).Err(); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
			} else {

				//send
			}
		case EventTypeLocation:
		case EventTypeTemplateSendJobFinish:
		}
	}
	//not same , denied
	c.JSON(http.StatusForbidden, "not accessed")
}

func sendTemplateMsg(toUser, url string) error {
	templateData := TemplateData{
		First:    DataItem{Value: "合适帐户创建成功"},
		Keyword1: DataItem{Value: "合适总部"},
		Keyword2: DataItem{Value: "刚刚"},
		Remark:   DataItem{Value: "点击这里进入我的账户 >>", Color: "#01934d"},
	}
	d, _ := json.Marshal(templateData)
	templateMessage := template.TemplateMessage{
		ToUser:     toUser,
		TemplateId: "lygCueaFhh-nXhu59WGoFzAfPZLOR2ZNJUbZeAYi8xE",
		URL:        url,
		Data:       []byte(d),
	}
	msgID, err := template.Send(wechatClient, templateMessage)
	if err != nil {
		return err
	}
	fmt.Println(msgID)
	return nil
}
