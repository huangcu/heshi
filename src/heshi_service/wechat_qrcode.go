package main

import (
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
