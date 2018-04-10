package main

import (
	"fmt"
	"heshi/errors"
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func parseExcel(c *gin.Context) {
	fheader, err := c.FormFile("upload")

	if err != nil {
		fmt.Println("here")
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	f, err := fheader.Open()
	if err != nil {
		fmt.Println("here1")
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		fmt.Println("here2")
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	rows := xlsx.GetRows("Sheet1")
	c.JSON(http.StatusOK, rows)
}
