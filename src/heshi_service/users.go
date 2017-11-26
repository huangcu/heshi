package main

import (
	"github.com/gin-gonic/gin"
)

type NewUserRequest struct {
	Username  string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
	Cellphone string `validate:"regexp=^\\+{0,1}0{0,1}62[0-9]+"`
	Email     string `validate:"nonzero"`
	Password  string `validate:"min=8"`
}

func newUser(c *gin.Context) {

}

func updateUser(c *gin.Context) {
	var usertype int
	switch usertype {
	case 0:
	case 1:
	case 2:
	}

}

func getUser(c *gin.Context) {

}

func getAllUsers(c *gin.Context) {

}
