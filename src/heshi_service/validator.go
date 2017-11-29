package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"matc.com/util"

// 	"github.com/asaskevich/govalidator"
// 	"github.com/gin-gonic/gin"
// )

// const (
// 	VEMSG_SHOULD_BE_JSON      = "should be JSON"
// 	VEMSG_SHOULD_NOT_BE_EMPTY = "should not be empty"
// )

// var (
// 	//TODO not support web app type for now
// 	SUPPORTED_APPTYPE     = []string{"app", "web"}
// 	SUPPORTED_API         = []int{18, 19, 21, 22, 23, 24, 25}
// 	SUPPORTED_SCREEN      = []string{"HVGA", "QVGA", "WQVGA400", "WQVGA432", "WSVGA", "WVGA800", "WVGA854", "WXGA720", "WXGA800"}
// 	VALID_REINSTALL_VALUE = []int{0, 1, 2}
// 	// `validate:"min=6,max=40,regexp=^[a-zA-Z]*$"`
// // `validate:"regexp=^\\+{0,1}0{0,1}62[0-9]+"`
// // `validate:"nonzero"`
// // `validate:"min=8"`min=8,max=20,regexp=^[a-zA-Z0-9_`!@#$%^&.?()-=+]_$
// )

// type ValidationErrors map[string]string

// func (ve ValidationErrors) Error() string {
// 	var buff bytes.Buffer
// 	for key, msg := range ve {
// 		buff.WriteString(fmt.Sprintf("%s %s\n", key, msg))
// 	}
// 	return strings.TrimSpace(buff.String())
// }

// /*preValidateNewJob validate new job Params*/
// func preValidateNewJob(c *gin.Context, jobType string) ValidationErrors {
// 	ve := ValidationErrors{}

// 	appType := c.PostForm("app_type")
// 	if !util.IsInArrayString(appType, SUPPORTED_APPTYPE) {
// 		ve["app_type"] = "should be app or web"
// 	}

// 	if strings.Contains(jobType, "device") {
// 		if !govalidator.IsUUIDv4(c.PostForm("device_id")) {
// 			ve["device_id"] = "should be uuid v4"
// 		}
// 	} else if appType == "wechat" {
// 		ve["app_type"] = "wechat is not supported on emulator"
// 	}

// 	for k, v := range map[string]string{
// 		"script_url":     c.PostForm("script_url"),
// 		"result_url":     c.PostForm("result_url"),
// 		"job_status_url": c.PostForm("job_status_url"),
// 	} {
// 		if !govalidator.IsURL(v) {
// 			ve[k] = "should be URL"
// 		}
// 	}
// 	if appType != "wechat" {
// 		if !govalidator.IsURL(c.PostForm("app_url")) {
// 			ve["app_url"] = "should be URL"
// 		}

// 		setting := c.PostForm("setting")
// 		if !govalidator.IsJSON(setting) {
// 			ve["setting"] = VEMSG_SHOULD_BE_JSON
// 		} else {
// 			var s Setting
// 			if err := json.Unmarshal([]byte(setting), &s); err != nil {
// 				ve["setting"] = fmt.Sprintf("invalid json: %s", setting)
// 			}
// 			if strings.Contains(jobType, "emulator") {
// 				if !util.IsInArrayInt(s.API, SUPPORTED_API) {
// 					ve["setting_api"] = fmt.Sprintf("API: %d is not supported.", s.API)
// 				}
// 				if !util.IsInArrayString(s.Screen, SUPPORTED_SCREEN) {
// 					ve["settings_screen"] = fmt.Sprintf("screen: %s is not supported.", s.Screen)
// 				}
// 			} else {
// 				if !util.IsInArrayInt(s.Reinstall, VALID_REINSTALL_VALUE) {
// 					ve["setting_reinstall"] = fmt.Sprintf("reinstall param: %d is not valid.", s.Reinstall)
// 				}
// 			}
// 		}
// 	}
// 	if len(ve) != 0 {
// 		return ve
// 	}
// 	return nil
// }
