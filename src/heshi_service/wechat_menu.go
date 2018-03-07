package main

import (
	"encoding/xml"
	"heshi/errors"
	"net/http"
	"util"

	"github.com/gin-gonic/gin"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
)

func createMenu(c *gin.Context) {
	b1 := menu.Button{
		Type: "click",
		Name: "请来问我",
		Key:  "KEY_KEFU",
	}
	b2 := menu.Button{
		Name: "商品选购",
	}
	b2s1 := menu.Button{
		Type: "view",
		Name: "本周特惠",
		URL:  "http://beyoudiamond.com/diamond-of-this-week.php?platform=weixin",
	}
	b2s2 := menu.Button{
		Type: "view",
		Name: "精选推荐钻石",
		URL:  "http://beyoudiamond.com/recommend-diamonds.php?platform=weixin",
	}
	b2s3 := menu.Button{
		Type: "view",
		Name: "挑选钻石",
		URL:  "http://beyoudiamond.com/diamonds.php?platform=weixin",
	}
	b2s4 := menu.Button{
		Type: "view",
		Name: "挑选首饰",
		URL:  "http://beyoudiamond.com/jewelry.PHP?platform=weixin",
	}
	b2s5 := menu.Button{
		Type: "view",
		Name: "挑选空托",
		URL:  "http://beyoudiamond.com/jewelry.php?class=mounting&platform=weixin",
	}
	b2.SubButtons = append(b2.SubButtons, b2s1, b2s2, b2s3, b2s4, b2s5)

	b3 := menu.Button{
		Name: "个人中心",
	}
	b3s1 := menu.Button{
		Type: "view",
		Name: "账户信息",
		URL:  "http://beyoudiamond.com/myaccount.php?platform=weixin",
	}
	b3s2 := menu.Button{
		Type: "view",
		Name: "我的积分",
		URL:  "http://beyoudiamond.com/myaccount.php?platform=weixin#section-mypoints",
	}
	// b3s3 := menu.Button{
	// 	Type:    "view",
	// 	Name:    "总部地址向导",
	// 	MediaId: "R575haACuEdixb-MfiD4pv-YAJC1eYjM1e5UC48WvnA",
	// }
	b3s4 := menu.Button{
		Type: "view",
		Name: "推荐给朋友",
		URL:  "http://beyoudiamond.com/recommendtofriend-weixin.php?platform=weixin",
	}
	b3.SubButtons = append(b3.SubButtons, b3s1, b3s2, b3s4)
	m := menu.Menu{
		Buttons: []menu.Button{b1, b2, b3},
	}

	//{"errcode":0,"errmsg":"ok"} create success
	if err := menu.Create(wechatClient, &m); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")
}

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141016
func handleMenuEvent(msg core.MixedMsg) {
	switch msg.EventType {
	case "CLICK":
		handleMenuClick(msg)
	case "VIEW":
	}
}

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140543
// 被动回复用户消息
// <xml>
// <ToUserName>< ![CDATA[toUser] ]></ToUserName>
// <FromUserName>< ![CDATA[fromUser] ]></FromUserName>
// <CreateTime>12345678</CreateTime>
// <MsgType>< ![CDATA[text] ]></MsgType>
// <Content>< ![CDATA[你好] ]></Content>
// </xml>
type replyMsgHeader struct {
	ToUserName   CDATAText `xml:"ToUserName"  `
	FromUserName CDATAText `xml:"FromUserName"`
	CreateTime   int64     `xml:"CreateTime"  `
	MsgType      CDATAText `xml:"MsgType"     `
}
type CDATAText struct {
	Text string `xml:",innerxml"`
}

type autoReplyMsg struct {
	XMLName struct{} `xml:"xml" json:"-"`
	replyMsgHeader
	Content CDATAText `xml:"Content" json:"Content"`
}

type transferToCustomerServiceReply struct {
	XMLName struct{} `xml:"xml" json:"-"`
	replyMsgHeader
	TransInfo *transInfo `xml:"TransInfo,omitempty" json:"TransInfo,omitempty"`
}

type transInfo struct {
	KfAccount CDATAText `xml:"KfAccount" json:"KfAccount"`
}

// <xml>
// <ToUserName><![CDATA[toUser]]></ToUserName>
// <FromUserName><![CDATA[FromUser]]></FromUserName>
// <CreateTime>123456789</CreateTime>
// <MsgType><![CDATA[event]]></MsgType>
// <Event><![CDATA[CLICK]]></Event>
// <EventKey><![CDATA[EVENTKEY]]></EventKey>
// </xml>
func handleMenuClick(msg core.MixedMsg) (string, error) {
	var reply string
	if msg.EventKey == "KEY_KEFU" {
		reply = "您好，有什么可以帮您的？"
	}
	art := autoReplyText(msg, reply)
	bs, err := xml.Marshal(art)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// <xml>
// <ToUserName><![CDATA[toUser]]></ToUserName>
// <FromUserName><![CDATA[FromUser]]></FromUserName>
// <CreateTime>123456789</CreateTime>
// <MsgType><![CDATA[event]]></MsgType>
// <Event><![CDATA[VIEW]]></Event>
// <EventKey><![CDATA[www.qq.com]]></EventKey>
// <MenuId>MENUID</MenuId>
// </xml>
func handleMenuView(msg core.MixedMsg) {
	//only to trace user activity
	util.Tracef("%s click %s", msg.FromUserName, msg.EventKey)
}
