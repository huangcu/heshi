package main

const (
	DIAMOND      = "diamond"
	SMALLDIAMOND = "small_diamond"
	JEWELRY      = "jewelry"
)

var VALID_PRODUCTS = []string{DIAMOND, SMALLDIAMOND, JEWELRY}
var diamondHeaders = []string{
	"diamond_id",
	"stock_ref",
	"shape",
	"carat",
	"color",
	"clarity",
	"grading_lab",
	"certificate_number",
	"cut_grade",
	"polish",
	"symmetry",
	"fluorescence_intensity",
	"country",
	"supplier",
	"price_no_added_value",
	"price_retail",
	"clarity_number",
	"cut_number",
}

var smallDiamondHeaders = []string{
	"size_from",
	"size_to",
	"price",
	"quantity",
}

//  "BR" :"圆形" /_images/constant/ico-stones.png
//  "PS": "梨形" /_images/constant/ico-stones.png
// 	"PR": "公主方" /_images/constant/ico-stones.png
// 	"HS": "心形" /_images/constant/ico-stones.png
// 	"MQ": "橄榄形" /_images/constant/ico-stones.png
// 	"OV": "椭圆形" /_images/constant/ico-stones.png
// 	"EM": "祖母绿" /_images/constant/ico-stones.png
// 	"RA": "雷蒂恩" /_images/constant/ico-stones.png???? not in db
// 	"CU": "枕形" /_images/constant/ico-stones.png
// 	"AS": "Asscher" /_images/constant/ico-stones.png
var VALID_DIAMOND_SHAPE = []string{
	"BR",
	"PS",
	"PR",
	"HS",
	"MQ",
	"OV",
	"EM",
	"CU",
	"AS",
	"RAD",
	"RBC",
	"RCRB",
	"CUSHION",
	"MARQUISE",
	"RC",
	"PE",
	"HT",
	"CMB",
}

//GIA: https://my.hrdantwerp.com/?L=&record_number='["certificate_number"].'&certificatetype=MC"
//HRD: http://www.gia.edu/cs/Satellite?pagename=GST%2FDispatcher&childpagename=GIA%2FPage%2FReportCheck&c=Page&cid=1355954554547&reportno='["certificate_number"].'"
//IGI: http://www.igiworldwide.com/verify.php?r='["certificate_number"].'"
var VALID_GRADING_LAB = []string{
	"GIA",
	"HRD",
	"IGI",
}

var VALID_CLARITY = []string{
	"IF",
	"VS1",
	"VS2",
	"VVS1",
	"VVS2",
	"SI1",
	"SI2",
	"I1",
	"I2",
	"P1",
}

var VALID_SYMMETRY = []string{
	"EX",
	"VG",
	"G",
	"GD",
	"Excellent",
	"EXC",
	"F",
	"FAIR",
}

var VALID_FLUORESCENCE_INTENSITY = []string{
	"None",
	"Medium",
	"Strong",
	"Faint",
	"NON",
	"Very Strong",
	"FNT",
	"NIL",
	"STR",
	"STG",
	"MED",
	"VST",
	"VSL",
	"SLT",
	"Slight",
	"Very Slight",
	"VSTG",
}

var VALID_POLISH = []string{
	"EX",
	"VG",
	"Excellent",
	"EXC",
	"G",
	"GD",
}

var VALID_SUPPLIER_NAME = []string{
	"KGK",
	"DIAM",
	"BEYOU-HESHI",
	"SUN",
	"HKEXPO",
	"HESHI",
}

var VALID_CUT_GRADE = []string{
	"EX",
	"G",
	"VG",
	"F",
	"GD",
	"Excellent",
	"EXC",
	"3EX",
	"NA",
	"NULL",
	"FAIR",
	"NN",
}

// var VALID_CLARITY_NUMBER = []string{
// 	"0", "1", "2", "3", "4", "5",
// }

var VALID_CUT_NUMBER = []string{
	"0", "1", "2", "3",
}

var VALID_COLOR = []string{}

// | J              |
// | G              |
// | F              |
// | H              |
// | I              |
// | E              |
// | M              |
// | L              |
// | D              |
// | K              |
// | L, Faint Brown |
// | FY             |
// | N              |
// | -              |
// | O-P            |
// | S-T            |
// | FLY            |
// | M, Faint Brown |
// |                |
// | ??             |
// | K, Faint Brown |
// | Y-Z            |
// | FLBY           |
// | ???            |
// | U-V            |
// | FBY            |
// | NFBY           |
// | NFLY           |
// | NFY            |
// | Q-R            |
// | NFLBY          |
// | U-Z            |
// | W-X            |
// | FLBGY          |
// | NBY            |
// | FANCY YELLOW   |
// | N-O            |
// | FDBY
