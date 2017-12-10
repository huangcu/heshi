package main

import (
	"log"
	"net/http"

	"gopkg.in/chanxuehong/wechat.v2/oauth2"

	"github.com/chanxuehong/rand"
	"github.com/gin-gonic/gin"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

// https://mp.weixin.qq.com
// CYHuang
// Lxxxxxx!!
// AppID wx7147ea39c1a30036
// AppSecret ba9e572ca65c6000c70cdb159254e32c
// 完成开发者设置，如果要成功调用access_token，你还需要设置IP白名单 (15.211.201.88/90)

func wechatAuth(c *gin.Context) {
	state := string(rand.NewHex())
	authURL := mpoauth2.AuthCodeURL(wxAppId, "https://localhost:8443/wechat/token", "snsapi_info", state)
	log.Println("AuthCodeURL:", authURL)
	c.Redirect(http.StatusFound, authURL)
}

//TODO not UI related should
func wechatToken(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Println("用户禁止授权")
		c.JSON(http.StatusOK, "用户禁止授权")
		return
	}

	queryState := c.Query("state")
	if queryState == "" {
		log.Println("state 参数为空")
		c.JSON(http.StatusOK, "state 参数为空")
		return
	}
	//TODO session, state
	// if savedState != queryState {
	// 	str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
	// 	io.WriteString(w, str)
	// 	log.Println(str)
	// 	return
	// }
	oauth2Client := oauth2.Client{
		Endpoint: endPoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		log.Println(err)
		return
	}
	//TODO write openid to db??
	log.Printf("token: %+v\r\n", token)
	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		log.Println(err)
		return
	}
	//TODO based on userinfo (openid, nickname, sex, city, province, country headimageurl, privilege,unionid)
	log.Printf("userinfo: %+v\r\n", userinfo)
}
