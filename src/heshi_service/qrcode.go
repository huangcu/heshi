package main

import (
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/mp/qrcode"
)

const (
	wxAppId     = "appid"
	wxAppSecret = "appsecret"

	wxOriId         = "oriid"
	wxToken         = "token"
	wxEncodedAESKey = "aeskey"
	wxOpenID        = "openid"
)

var (
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
	wechatClient      *core.Client           = core.NewClient(accessTokenServer, nil)
)

func getAccessToken() (string, error) {
	token, err := accessTokenServer.Token()
	if err != nil {
		return "", nil
	}
	return token, nil
}

func wechatLogin() {

	oauth2.AuthCodeURL(wxAppId, "http://heshi.com", "", "")
	oauth2.Auth(wxToken, wxOpenID, nil)
}

func qrCodePic() (string, error) {
	//sceneID
	tempQRCode, err := qrcode.CreateTempQrcode(wechatClient, 10000, 60)
	if err != nil {
		return "", err
	}

	return qrcode.QrcodePicURL(tempQRCode.Ticket), nil
}
