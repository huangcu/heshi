package main

import (
	"fmt"

	"gopkg.in/chanxuehong/wechat.v2/mp/dkf"
	kfaccount "gopkg.in/chanxuehong/wechat.v2/mp/dkf/account"
	"gopkg.in/chanxuehong/wechat.v2/mp/dkf/session"
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
func sendKfMsg(toUser, content, kfAccount string) error {
	t := custom.Text{
		MsgHeader: custom.MsgHeader{
			ToUser:  toUser,
			MsgType: custom.MsgTypeText,
		},
		//Creating anonymous structures
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
		CustomService: &custom.CustomService{
			KfAccount: kfAccount,
		},
	}
	if err := logTextMsgToDB(toUser, kfAccount, content, "TO"); err != nil {
		return err
	}
	return custom.Send(wechatClient, t)
}

func createSession(fromUser, kfAccount, content string) error {
	return session.Create(wechatClient, fromUser, kfAccount, content)
}

func closeSession(fromUser, kfAccount, content string) error {
	return session.Close(wechatClient, fromUser, kfAccount, content)
}

func waitCaseNumber() (int, error) {
	waitCaseList, err := session.WaitCaseList(wechatClient)
	if err != nil {
		return 0, err
	}
	return waitCaseList.TotalCount, nil
}

func isKfAvaiable(kfAccount string) (bool, error) {
	kfList, err := dkf.OnlineKfList(wechatClient)
	if err != nil {
		return false, err
	}

	for _, kf := range kfList {
		if kfAccount == kf.Account {
			//kf online and accepting any case is not exceed auto accept number
			if (kf.Status == 1 || kf.Status == 2 || kf.Status == 3) && kf.AcceptingNumber < kf.AutoAcceptNumber {
				return true, nil
			}
		}
	}
	return false, nil
}

func logTextMsgToDB(user, kfAccout, content, direction string) error {
	q := fmt.Sprintf(`INSERT into messages (id, user, content, kf_account, direction) 
	VALUES('%s','%s','%s','%s','%s')`,
		newV4(), user, kfAccout, content, direction)
	_, err := dbExec(q)
	return err
}
