package main

import (
	"fmt"
	"heshi/errors"
	"time"
	"util"

	"gopkg.in/chanxuehong/wechat.v2/mp/message/custom"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140453
func handleTextMsg(msg core.MixedMsg) *replyText {
	q := `SELECT COUNT(*) FROM users WHERE wechat_openid=`
	var count int
	if err := dbQueryRow(q, msg.FromUserName).Scan(&count); err != nil {
		util.Traceln(errors.GetMessage(err))
	}
	if count != 1 {
		var reply string
		if msg.EventKey == "KEY_KEFU" {
			reply = "您尚未注册合适账户。建议您先创建一个合适账户，这样我们可以更好的为您服务。"
		}
		return &replyText{
			ToUserName:   CDATAText{Text: cdataStartLiteral + msg.FromUserName + cdataEndLiteral},
			FromUserName: CDATAText{Text: cdataStartLiteral + msg.ToUserName + cdataEndLiteral},
			MsgType:      CDATAText{Text: cdataStartLiteral + "Text" + cdataEndLiteral},
			CreateTime:   CDATAText{Text: fmt.Sprintf("%s%d%s", cdataStartLiteral, time.Now().Unix(), cdataEndLiteral)},
			Content:      CDATAText{Text: cdataStartLiteral + reply + cdataEndLiteral},
		}
	}
	//kf reply
	kfReply := custom.Text{
		MsgHeader: custom.MsgHeader{
			ToUser:  msg.FromUserName,
			MsgType: "Text",
		},
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "",
		},
		CustomService: &custom.CustomService{
			KfAccount: "test1@kftest",
		},
	}
	custom.Send(wechatClient, kfReply)

	return nil
}

func handleImageMsg(msg core.MixedMsg) {

}
