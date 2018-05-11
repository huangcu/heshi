package main

import (
	"fmt"
	"heshi/errors"
	"net/http"
	"util"

	"gopkg.in/chanxuehong/wechat.v2/oauth2"

	"github.com/chanxuehong/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

// https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
// gh_5bd700510a86
// appID wxa6c9fc631124397a
// appsecret ad23b9ed5679d5be74f69db6875dcd7f
// https://open.weixin.qq.com
// huang372923@sina.com
// Lxxxxxxx!!
// https://mp.weixin.qq.com
// CYHuang
// gh_9b9257df583f
// h******37***3@hotmail.com
// Lxxxxxx!!
// AppID wx7147ea39c1a30036
// AppSecret ba9e572ca65c6000c70cdb159254e32c
// 完成开发者设置，如果要成功调用access_token，你还需要设置IP白名单 (15.211.201.88/90/15.107.17.33)
const (
	// 	wxAppID          = "wx7147ea39c1a30036"               //wx02b69905483c2df2
	// 	wxAppSecret      = "ba9e572ca65c6000c70cdb159254e32c" //wx02b69905483c2df2&secret=029248abec380aaab05b95edf58681bd
	wxAppAccount     = "gh_5bd700510a86"
	wxAppIDDebug     = "wxa6c9fc631124397a"
	wxAppSecretDebug = "ad23b9ed5679d5be74f69db6875dcd7f"

	wxOriID           = "oriid"
	wxToken           = "token"
	wxEncodedAESKey   = "aeskey"
	cdataStartLiteral = "<![CDATA["
	cdataEndLiteral   = "]]>"
)

var (
	serverURI                                = "http://cee8c7f5.ngrok.io"
	redirectURI                              = serverURI + "/api/wechat/token"
	redirectLogin                            = serverURI + "/webpage/login.html"
	endPoint                                 = mpoauth2.NewEndpoint(wxAppIDDebug, wxAppSecretDebug) //*mpoauth2.Endpoint
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppIDDebug, wxAppSecretDebug, nil)
	wechatClient                             = core.NewClient(accessTokenServer, nil) //*core.Client
)

// ec2-52-221-233-143.ap-southeast-1.compute.amazonaws.com
// https://open.weixin.qq.com/connect/outh2/authorize?appid=wxa6c9fc631124397a&redirect_uri=http://cee8c7f5.ngrok.io/api/wechat/token&response_type=code&scope=snsapi_base#wechat_redirect
// $auth_url='https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx02b69905483c2df2
// &redirect_uri='.urlencode("http://www.beyoudiamond.com/authreceiver.php").
// '&response_type=code&scope=snsapi_base#wechat_redirect';
// 发起授权
// 微信扫码支付：此公众号并没有这些scope的权限，错误码：10005->
// scope 应用授权作用域。
// snsapi_base：静默授权，可获取成员的基础信息；不弹出授权页面，直接跳转，只能获取用户openid
// snsapi_userinfo：静默授权，可获取成员的详细信息，但不包含手机、邮箱；弹出授权页面，可通过openid拿到昵称、性别、所在地。并且，即使在未关注的情况下，只要用户授权，也能获取其信息
// snsapi_privateinfo：手动授权，可获取成员的详细信息，包含手机、邮箱。
func wechatAuth(c *gin.Context) {
	state := string(rand.NewHex())
	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, state)
	s.Save()
	authURL := mpoauth2.AuthCodeURL(wxAppIDDebug, redirectURI, "snsapi_userinfo", state)
	util.Traceln("AuthCodeURL:", authURL)
	c.Redirect(http.StatusFound, authURL)
}

//TODO not UI related should
func wechatToken(c *gin.Context) {
	code := c.Query("code")
	util.Println("code", code)
	if code == "" {
		util.Println("用户禁止授权")
		c.Redirect(http.StatusFound, redirectLogin)
		return
	}
	queryState := c.Query("state")
	util.Println("state", queryState)
	if queryState == "" {
		util.Println("state 参数为空")
		c.Redirect(http.StatusFound, redirectLogin)
		return
	}

	s := sessions.Default(c)
	savedState := s.Get(USER_SESSION_KEY)
	if savedState != queryState {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
		util.Println(str)
		c.Redirect(http.StatusFound, redirectLogin)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: endPoint,
	}
	util.Println("get access token")
	// https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code
	token, err := oauth2Client.ExchangeToken(code)

	if err != nil {
		util.Println("error to get access token", err)
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	util.Printf("token: %+v\r\n", token)
	exist, err := isWechatUserExist(token.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	if !exist {
		// https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
		userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}

		wu := wechatUserInfo{userinfo}
		if err := wu.newWechatUser(); err != nil {
			c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
			return
		}
		util.Printf("userinfo: %+v\r\n", userinfo)
	}
	s.Set(USER_SESSION_KEY, token.OpenId)
	s.Save()
	util.Println("redirect to home page: http://localhost:8081")
	c.Redirect(http.StatusFound, "http://localhost:8081")
}

// subscribe	用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
// openid	用户的标识，对当前公众号唯一
// nickname	用户的昵称
// sex	用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
// city	用户所在城市
// country	用户所在国家
// province	用户所在省份
// language	用户的语言，简体中文为zh_CN
// headimgurl	用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
// subscribe_time	用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
// unionid	只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
// remark	公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
// groupid	用户所在的分组ID（兼容旧的用户分组接口）
// tagid_list	用户被打上的标签ID列表
