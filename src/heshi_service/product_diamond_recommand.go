package main

import (
	"fmt"
	"util"

	"github.com/gin-gonic/gin"
)

//TODO another field to indicate recommand by heshi
func recommnadDiamonds(c *gin.Context) {
	///api/recommand/diamonds?sort='carat'&order='up'
	sort := "carat"
	if c.Query("sort") == "price_retail" {
		sort = "price_retail"
	}
	order := "ASC"
	if c.Query("order") == "down" {
		order = "DESC"
	}
	q := fmt.Sprintf("SELECT * FROM diamonds WHERE recommnaded_by_heshi=='YES' AND statu <> 'SOLD' ORDER BY %s %s",
		sort, order)
	util.Println(q)
	// $query_sort=' ORDER BY carat ';
	// $query_order=' ASC';
	// if(isset($_GET['sort'])){
	// 	if($_GET['sort']=='c'){
	// 		$query_sort=' ORDER BY carat ';
	// 	}else if($_GET['sort']=='p'){
	// 		$query_sort=' ORDER BY price_retail ';
	// 	}
	// }
	// if(isset($_GET['order'])){
	// 	if($_GET['order']=='up'){
	// 		$query_order=' ASC';
	// 	}else if($_GET['order']=='down'){
	// 		$query_order=' DESC';
	// 	}
	// }
	// SELECT * FROM diamonds WHERE pic1 IS NOT NULL AND status <> "SOLD" '.$query_sort.$query_order;
}
