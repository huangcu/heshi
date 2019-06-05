package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"heshi/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"util"

	"gopkg.in/chanxuehong/wechat.v2/mp/dkf/session"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140453
func handleTextMsg(msg core.MixedMsg) (string, error) {
	ss, err := session.Get(wechatClient, msg.FromUserName)
	if err != nil {
		return "", err
	}
	if err := logTextMsgToDB(msg.FromUserName, ss.KfAccount, msg.Content, "FROM"); err != nil {
		util.Traceln(errors.GetMessage(err))
	}
	//already in a session
	//用户被客服接入以后，客服关闭会话以前，处于会话过程中时，用户发送的消息均会被直接转发至客服系统, return directly
	if ss.KfAccount != "" {
		return "", nil
	}

	//if not in a session
	q := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE wechat_openid='%s'`, msg.FromUserName)
	var count int
	if err := dbQueryRow(q).Scan(&count); err != nil {
		util.Traceln(errors.GetMessage(err))
		return "", err
	}
	if count != 1 {
		reply := "您尚未注册合适账户。建议您先创建一个合适账户，这样我们可以更好的为您服务。"
		art := autoReplyText(msg, reply)
		bs, err := xml.Marshal(art)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}

	//transfer to kf
	q = fmt.Sprintf("SELECT kf_account FROM wechat_messages WHERE user='%s' ORDER BY created_at DESC", msg.FromUserName)
	var kfAccount string
	if err := dbQueryRow(q).Scan(&kfAccount); err != nil {
		//send to 多客服
		if err == sql.ErrNoRows {
			art := passToAllKf(msg)
			bs, err := xml.Marshal(art)
			if err != nil {
				return "", err
			}
			return string(bs), nil
		}
		util.Traceln(errors.GetMessage(err))
		return "", err
	}

	//select the latest kf served, and check if kf available
	available, err := isKfAvailable(kfAccount)
	if err != nil {
		return "", err
	}
	if !available {
		art := passToAllKf(msg)
		bs, err := xml.Marshal(art)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}

	art := passToKf(msg, kfAccount)
	bs, err := xml.Marshal(art)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// <xml><ToUserName><![CDATA[gh_5bd700510a86]]></ToUserName>
// <FromUserName><![CDATA[om7Sh0xuvSL4h4RSzbwu1qbwePtw]]></FromUserName>
// <CreateTime>1520414204</CreateTime>
// <MsgType><![CDATA[image]]></MsgType>
// <PicUrl><![CDATA[http://mmbiz.qpic.cn/mmbiz_jpg/IRxsVrZTsNAxspBTFlZH0f11wDHy1F1Y0ya8eV4zmLLNV2tnPNO7BicpdXS4icicr5yekUTpUoFRmhZA7kw9w2eiaA/0]]></PicUrl>
// <MsgId>6530129282993985415</MsgId>
// <MediaId><![CDATA[mCVJIBFT2-SDK3jZ9VyqwfMocdsD2eq1aSrauP1mODGZ8ZUaDikCRZU2wLyz23Qk]]></MediaId>
func handleImageMsg(msg core.MixedMsg) error {
	ss, err := session.Get(wechatClient, msg.FromUserName)
	if err != nil {
		return err
	}
	//Download pic, save to filepath.Join(".wechat", FromUser, .time.now().unix().JPG)
	response, err := http.Get(msg.PicURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := os.MkdirAll(filepath.Join(".wechat", msg.FromUserName), 0755); err != nil {
		return err
	}
	//open a file for writing
	picURL := filepath.Join(".wechat", msg.FromUserName, fmt.Sprintf("%d.jpg", time.Now().Unix()))
	file, err := os.Create(picURL)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	file.Close()

	if err := logImageMsgToDB(msg.FromUserName, ss.KfAccount, picURL, "FROM"); err != nil {
		util.Traceln(errors.GetMessage(err))
	}
	//already in a session
	//用户被客服接入以后，客服关闭会话以前，处于会话过程中时，用户发送的消息均会被直接转发至客服系统, return directly
	if ss.KfAccount != "" {
		return nil
	}

	return nil
}

func passToAllKf(msg core.MixedMsg) *transferToCustomerServiceReply {
	return &transferToCustomerServiceReply{
		replyMsgHeader: replyMsgHeader{
			ToUserName:   CDATAText{Text: cdataStartLiteral + msg.FromUserName + cdataEndLiteral},
			FromUserName: CDATAText{Text: cdataStartLiteral + msg.ToUserName + cdataEndLiteral},
			MsgType:      CDATAText{Text: cdataStartLiteral + string(response.MsgTypeTransferCustomerService) + cdataEndLiteral},
			CreateTime:   time.Now().Unix(),
		},
	}
}

func passToKf(msg core.MixedMsg, kfAccount string) *transferToCustomerServiceReply {
	return &transferToCustomerServiceReply{
		replyMsgHeader: replyMsgHeader{
			ToUserName:   CDATAText{Text: cdataStartLiteral + msg.FromUserName + cdataEndLiteral},
			FromUserName: CDATAText{Text: cdataStartLiteral + msg.ToUserName + cdataEndLiteral},
			MsgType:      CDATAText{Text: cdataStartLiteral + string(response.MsgTypeTransferCustomerService) + cdataEndLiteral},
			CreateTime:   time.Now().Unix(),
		},
		TransInfo: &transInfo{
			KfAccount: CDATAText{Text: cdataStartLiteral + kfAccount + cdataEndLiteral},
		},
	}
}

func autoReplyText(msg core.MixedMsg, content string) *autoReplyMsg {
	return &autoReplyMsg{
		replyMsgHeader: replyMsgHeader{
			ToUserName:   CDATAText{Text: cdataStartLiteral + msg.FromUserName + cdataEndLiteral},
			FromUserName: CDATAText{Text: cdataStartLiteral + msg.ToUserName + cdataEndLiteral},
			MsgType:      CDATAText{Text: cdataStartLiteral + string(response.MsgTypeText) + cdataEndLiteral},
			CreateTime:   time.Now().Unix(),
		},
		Content: CDATAText{Text: cdataStartLiteral + content + cdataEndLiteral},
	}
}
