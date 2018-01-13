package errors

import (
	"fmt"
)

type HSMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func showWithCode(code int, err error) HSMessage {
	return HSMessage{code, GetMessage(err)}
}

func Show(code int, err interface{}) HSMessage {
	if msg, ok := err.(string); ok {
		return HSMessage{code, msg}
	}
	e := err.(error)
	println(e.Error())
	return showWithCode(code, e)
}

func Showf(code int, format string, a ...interface{}) HSMessage {
	s := fmt.Sprintf(format, a...)
	print(s)
	return Show(code, s)
}
