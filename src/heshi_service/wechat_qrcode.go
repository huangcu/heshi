package main

import (
	"math"
	"net/http"
	"util"

	"github.com/chanxuehong/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	QRCODE "github.com/skip2/go-qrcode"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/mp/qrcode"
)

// ①用户登录网页，点击“绑定微信账户”；
// ②后台使用微信接口，生成二维码链接返回给前端显示，并建立场景值A与用户的对应关系；
// ③用户扫描二维码，并点击关注微信公众号（假如已关注，直接跳到④）；
// ④后台接收微信服务器推送的场景值A；
// ⑤后台根据场景值A，查询到对应的用户ID（依赖于②中建立的对应关系）；
// ⑥建立用户userid与微信用户openid的对应关系；
// ⑦给用户的微信客户端推送“绑定成功”的提示；
// ⑧通知前台页面，绑定已完成，刷新页面，并返回一些微信账户信息。完成绑定
var sceneID int32

type TempQrCode struct {
	SceneID   int32  `json:"scene_id"`
	QrCodeURL string `json:"qr_code_url"`
}

func wechatQrCode(c *gin.Context) {
	state := string(rand.NewHex())
	s := sessions.Default(c)
	s.Set(USER_SESSION_KEY, state)
	s.Save()
	authURL := mpoauth2.AuthCodeURL(wxAppIDDebug, redirectURI, "snsapi_userinfo", state)
	util.Println("qrcode AuthCodeURL:", authURL)
	QRCODE.WriteFile(authURL, QRCODE.Medium, 256, "qr.png")
}

func wechatTempQrCode(c *gin.Context) {
	// 临时二维码的scene_id为32位非0整型->是32位的二进制数，即最大值是2的32次方减1也就是4294967295
	if sceneID > math.MaxInt32 || sceneID == 0 {
		sceneID = sceneID + 1
	} else {
		sceneID = 1
	}

	tempQRCode, err := qrcode.CreateTempQrcode(wechatClient, sceneID, 120)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	qrcodeURL := qrcode.QrcodePicURL(tempQRCode.Ticket)
	t := TempQrCode{
		SceneID:   sceneID,
		QrCodeURL: qrcodeURL,
	}
	util.Printf("%v", t)
	c.JSON(http.StatusOK, t)
}

func wechatQrCodeStatus(c *gin.Context) {
	sceneID := c.PostForm("scene_id")
	openID, err := redisClient.Get(sceneID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, "")
		return
	}
	if _, err := redisClient.Del(sceneID).Result(); err != nil {
		util.Printf("fail to clean key %s from redis db", sceneID)
	}
	c.JSON(http.StatusOK, openID)
}
