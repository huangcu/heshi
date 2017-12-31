package main

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

//微信公众号平台Callback - 接入验证
func wechatCallback(c *gin.Context) {
	log.Println("entry: wechatCallback")

	// log.Printf("request: %#v", c.Request)
	//接入验证
	echostr := c.Query("echostr")
	if echostr != "" {
		verified := checkSignature(c.Query("signature"), c.Query("timestamp"), c.Query("nonce"))
		if verified {
			//accessed
			c.String(http.StatusOK, echostr)
			return
		}
		c.String(http.StatusOK, "verify failed")
		return
	}
	// if (c.Request.Method == "POST")
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var msg MixedMsg
	log.Printf("msg body received from wechat server: %s", string(bs))
	err = xml.Unmarshal(bs, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("response msg from wechat %v", msg)
	switch msg.MsgType {
	case "event":
		wechatEventHandler(msg)
	case "text":
	case "image":
	case "voice":
	case "video", "shortvideo":
	case "location":
	case "link":
	}

	//TODO
	// 	假如服务器无法保证在五秒内处理并回复，必须做出下述回复，这样微信服务器才不会对此作任何处理，并且不会发起重试（这种情况下，可以使用客服消息接口进行异步回复），否则，将出现严重的错误提示。详见下面说明：
	// （推荐方式）直接回复success
	// 直接回复空串（指字节长度为0的空字符串，而不是XML结构体中content字段的内容为空）
	c.String(http.StatusOK, "")
}

func checkSignature(signature, timestamp, nonce string) bool {
	token := "BEYOU_SIHUI"
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)
	h := sha1.New()
	_, err := h.Write([]byte(strings.Join(tmpArr, "")))
	if err != nil {
		log.Printf("check signature failed. error: %s", err.Error())
		return false
	}
	bs := h.Sum(nil)
	s := fmt.Sprintf("%x", bs)

	if s == signature {
		//accessed
		return true
	}
	return false
}

func wechatEventHandler(msg MixedMsg) error {
	switch msg.EventType {
	case EventTypeSubscribe, EventTypeScan:
		if err := redisClient.Set(msg.EventKey, msg.FromUserName, 0).Err(); err != nil {
			log.Printf("fail to write to redis db. err: %s", err.Error())
			return err
		}
		sendTemplateMsg(msg.FromUserName, "http://721e2175.ngrok.io/api/wechat/auth")
	case EventTypeUnsubscribe:
	case EventTypeLocation:
	case EventTypeTemplateSendJobFinish:
		if msg.Status == "success" {
			return nil
		}
	}
	return nil
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
	log.Printf("message: %v sent out. msgid: %d", templateMessage, msgID)
	return nil
}
