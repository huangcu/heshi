package main

import (
	"fmt"

	"github.com/chanxuehong/rand"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/mp/qrcode"
)

const (
	wxAppId     = "wx7147ea39c1a30036"               //wx02b69905483c2df2
	wxAppSecret = "ba9e572ca65c6000c70cdb159254e32c" //wx02b69905483c2df2&secret=029248abec380aaab05b95edf58681bd

	wxOriId         = "oriid"
	wxToken         = "token"
	wxEncodedAESKey = "aeskey"
	wxOpenID        = "openid"
)

var (
	endPoint          *oauth2.Endpoint       = oauth2.NewEndpoint(wxAppId, wxAppSecret)
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

// $auth_url='https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx02b69905483c2df2
// &redirect_uri='.urlencode("http://www.beyoudiamond.com/authreceiver.php").
// '&response_type=code&scope=snsapi_base#wechat_redirect';
func getWechatOAuth20Code() (string, error) {
	// scope 应用授权作用域。
	// snsapi_base：静默授权，可获取成员的基础信息；不弹出授权页面，直接跳转，只能获取用户openid
	// snsapi_userinfo：静默授权，可获取成员的详细信息，但不包含手机、邮箱；弹出授权页面，可通过openid拿到昵称、性别、所在地。并且，即使在未关注的情况下，只要用户授权，也能获取其信息
	// snsapi_privateinfo：手动授权，可获取成员的详细信息，包含手机、邮箱。
	state := string(rand.NewHex())
	authURL := oauth2.AuthCodeURL(wxAppId, "http://heshi.com", "snsapi_info", state)
	fmt.Println(authURL)
	return "code", nil
}

// https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET
// &code=CODE&grant_type=authorization_code
func getWechatAccessToken(code string) (string, error) {
	exchangeTokenURL := fmt.Sprint("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wxAppId, wxAppSecret, code)
	fmt.Println(exchangeTokenURL)
	//write open id to DB
	return "", nil
}

// https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
func getWechatUserInfo(accessToken, openId string) {

}

func qrCodePic() (string, error) {
	//sceneID
	tempQRCode, err := qrcode.CreateTempQrcode(wechatClient, 10000, 60)
	if err != nil {
		return "", err
	}

	return qrcode.QrcodePicURL(tempQRCode.Ticket), nil
}
func end() {
	endPoint.ExchangeTokenURL("")
}
