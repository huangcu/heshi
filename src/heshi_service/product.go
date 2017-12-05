package main

import "github.com/gin-gonic/gin"

type product struct {
	Diamond      []diamond      `json:"diamond"`
	Jewelry      []jewelry      `json:"jewelry"`
	SmallDiamond []smallDiamond `json:"small_diamond"`
}

func getAllProducts(c *gin.Context) {

}
