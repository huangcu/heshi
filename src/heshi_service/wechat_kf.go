package main

import (
	kfaccount "gopkg.in/chanxuehong/wechat.v2/mp/dkf/account"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/custom"
)

func addKfAccount() error {
	return kfaccount.Add(wechatClient, "account", "nickname", "password", true)
}

func updateKfAccountHeadImg() error {
	return kfaccount.UploadHeadImage(wechatClient, "kfaccount", "imagefilepath")
}

// {
//     "touser":"OPENID",
//     "msgtype":"text",
//     "text":
//     {
//          "content":"Hello World"
//     }
// "customservice":
// {
//      "kf_account": "test1@kftest"
// }
// }
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140547
func kfSendMsg() {
	t := custom.Text{
		MsgHeader: custom.MsgHeader{
			ToUser:  "touser",
			MsgType: "text",
		},
		//Creating anonymous structures
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "",
		},
		CustomService: &custom.CustomService{
			KfAccount: "test1@kftest",
		},
	}
	custom.Send(wechatClient, t)
}
