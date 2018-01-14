package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//TODO search
func searchProducts(c *gin.Context) {
	category := c.Param("category")
	if category != "diamonds" || category != "jewelrys" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if category == "diamonds" {

		searchDiamonds()
	} else {
		searchJewelrys()
	}
}

func searchDiamonds() {

}

func searchJewelrys() {

}
