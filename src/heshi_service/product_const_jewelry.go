package main

var jewelryHeaders = []string{
	//must have fields
	"stock_id",
	"name",
	"need_diamond",
	"category",
	"mounting_type",
	"material",
	"metal_weight",
	"dia_shape",

	//optional fields
	"unit_number",
	"dia_size_min",
	"dia_size_max",
	"main_dia_num",
	"main_dia_size",
	"small_dias",
	"small_dia_num",
	"small_dia_carat",
	"price",
	"video_link",
	"text",
	"status",
	"verified",
	"featured",
	"stock_quantity",
	"profitable",
	"free_acc",
	"image",
	"image1",
	"image2",
	"image3",
	"image4",
	"image5",
}

var gemHeaders = []string{
	"name",
	"stock_id",
	"size",
	"material",
	"shape",
	"certificate",

	"price",
	"text",
	"status",
	"verified",
	"featured",
	"stock_quantity",
	"profitable",
	"free_acc",
	"image",
	"image1",
	"image2",
	"image3",
	"image4",
	"image5",
}

//TODO ref
// <option value="JP">素金吊坠／项链</option> 1
// <option value="JR">素金戒指</option> 2
// <option value="JE">素金耳环／耳钉</option> 3
// <option value="ZP">镶碎钻吊坠／项链</option> 1 | 5
// <option value="ZR">镶碎钻戒指</option> 2
// <option value="ZE">镶碎钻耳环／耳钉</option> 3

// <option value="CP">成品吊坠／项链</option> 1 | 5 /NO
// <option value="CR">成品戒指</option> 2 /NO
// <option value="CE">成品耳环／耳钉</option> 3/NO

// JP small_dias ="NO" AND need_diamond = "YES" AND (category = 1 OR category = 5)
// JR small_dias ="NO" AND need_diamond = "YES" AND category = 2
// JE small_dias ="NO" AND need_diamond = "YES" AND category = 3
// ZP small_dias ="YES" AND need_diamond = "YES" AND (category = 1 OR category = 5)
// ZR small_dias ="YES" AND need_diamond = "YES" AND category = 2
// ZE small_dias ="YES" AND need_diamond = "YES" AND category = 3
// CP small_dias = "NO" AND need_diamond = "NO" AND (category = 1 OR category = 5)
// CR small_dias = "NO" AND need_diamond = "NO" AND category = 2
// CE small_dias = "NO" AND need_diamond = "NO" AND category = 3

// jewelry_category
// id | category_en  | category_cn   | main_diamond_num
// |  1 | pendant      | 吊坠          |                1 |
// |  2 | ring         | 戒指          |                1 |
// |  3 | earring      | 耳环&耳钉     |                2 |
// |  9 | bracelet     | 手链          |                1 |
// |  5 | necklace     | 项链          |                1 |
// | 10 | precious-gem | 彩宝          |                1 |
//成品首饰 need diamond: NO,
//空托    need diamond: YES,

var VALID_CATEGORY = []string{
	"PENDANT",
	"RING",
	"EARRING",
	"BRACELET",
	"NECKLACE",
}

// Request URL:http://www.beyoudiamond.com/jewelry.php?class=mounting (kongtuo)
// Request URL:http://www.beyoudiamond.com/jewelry.php (chengpin)
// Request URL:http://www.beyoudiamond.com/colored-gems.php
var VALID_MATERIAL = []string{
	"PT",
	"ROSE_GOLD",
	"COLORED_GOLD",
	"UNKNOWN",
}

var VALID_MOUNTING_TYPE = []string{
	"3NODE",
	"4NODE",
	"6NODE",
	"SURROUND",
	"SPECIAL",
}

//dia_shape, should be array;
