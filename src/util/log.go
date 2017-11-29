package util

import (
	"fmt"
	"log"
	"strings"
)

var Logger *log.Logger

func Println(a ...interface{}) {
	if Logger == nil {
		fmt.Println(a...)
		return
	}
	Logger.Println(a...)
}

func Printf(format string, a ...interface{}) {
	if Logger == nil {
		fmt.Printf(format, a...)
		return
	}
	Logger.Printf(format, a...)
}

func PrintCmd(cmd string, ps ...string) {
	Println(getCmd(cmd, ps...))
}

func Tracef(format string, a ...interface{}) {
	if ShouldTrace() {
		Printf(format, a...)
	}
}

func Traceln(a ...interface{}) {
	if ShouldTrace() {
		Println(a...)
	}
}

func TraceError(err error) {
	if err != nil && ShouldTrace() {
		// Println(errors.GetMessage(err))
	}
}

func TraceCmd(cmd string, ps ...string) {
	if ShouldTrace() {
		Println("\t->", getCmd(cmd, ps...))
	}
}

func getCmd(cmd string, ps ...string) string {
	return fmt.Sprintf("%s %s", cmd,
		strings.TrimSpace(strings.Trim(fmt.Sprint(ps), "[]")))
}
