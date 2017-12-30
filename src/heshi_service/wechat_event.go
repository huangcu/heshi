package main

import "gopkg.in/chanxuehong/wechat.v2/mp/core"

const (
	// 普通事件类型
	EventTypeSubscribe             EventType = "subscribe"   // 关注事件, 包括点击关注和扫描二维码(公众号二维码和公众号带参数二维码)关注
	EventTypeUnsubscribe           EventType = "unsubscribe" // 取消关注事件
	EventTypeScan                  EventType = "SCAN"        // 已经关注的用户扫描带参数二维码事件
	EventTypeLocation              EventType = "LOCATION"    // 上报地理位置事件
	EventTypeTemplateSendJobFinish EventType = "TEMPLATESENDJOBFINISH"
)

const (
	TemplateSendStatusSuccess            = "success"               // 送达成功时
	TemplateSendStatusFailedUserBlock    = "failed:user block"     // 送达由于用户拒收(用户设置拒绝接收公众号消息)而失败
	TemplateSendStatusFailedSystemFailed = "failed: system failed" // 送达由于其他原因失败
)

type TemplateSendJobFinishEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	core.MsgHeader
	EventType EventType `xml:"Event"  json:"Event"`  // 此处为 TEMPLATESENDJOBFINISH
	MsgId     int64     `xml:"MsgId"  json:"MsgId"`  // 模板消息ID
	Status    string    `xml:"Status" json:"Status"` // 发送状态
}

// 关注事件
type SubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MsgHeader
	EventType EventType `xml:"Event" json:"Event"` // subscribe

	// 下面两个字段只有在扫描带参数二维码进行关注时才有值, 否则为空值!
	EventKey string `xml:"EventKey,omitempty" json:"EventKey,omitempty"` // 事件KEY值, 格式为: qrscene_二维码的参数值
	Ticket   string `xml:"Ticket,omitempty"   json:"Ticket,omitempty"`   // 二维码的ticket, 可用来换取二维码图片
}

// 取消关注事件
type UnsubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MsgHeader
	EventType EventType `xml:"Event"              json:"Event"`              // unsubscribe
	EventKey  string    `xml:"EventKey,omitempty" json:"EventKey,omitempty"` // 事件KEY值, 空值
}

// 用户已关注时, 扫描带参数二维码的事件
type ScanEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MsgHeader
	EventType EventType `xml:"Event"    json:"Event"`    // SCAN
	EventKey  string    `xml:"EventKey" json:"EventKey"` // 事件KEY值, 二维码的参数值(scene_id, scene_str)
	Ticket    string    `xml:"Ticket"   json:"Ticket"`   // 二维码的ticket, 可用来换取二维码图片
}

// 上报地理位置事件
type LocationEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	MsgHeader
	EventType EventType `xml:"Event"     json:"Event"`     // LOCATION
	Latitude  float64   `xml:"Latitude"  json:"Latitude"`  // 地理位置纬度
	Longitude float64   `xml:"Longitude" json:"Longitude"` // 地理位置经度
	Precision float64   `xml:"Precision" json:"Precision"` // 地理位置精度(整数? 但是微信推送过来是浮点数形式)
}

type TemplateData struct {
	First    DataItem `json:"first"`
	Keyword1 DataItem `json:"keyword1"`
	Keyword2 DataItem `json:"keyword2"`
	Remark   DataItem `json:"remark"`
}

type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}
