package main

const (
	DIAMOND      = "diamond"
	SMALLDIAMOND = "small_diamond"
	JEWELRY      = "jewelry"
	GEM          = "gem"
)

var VALID_PRODUCTS = []string{DIAMOND, SMALLDIAMOND, JEWELRY}

var diamondHeaders = []string{
	//must have fields
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

	// optional fields
	"price_retail",
	"featured",
	"recommand_words",
	"extra_words",
	"image",
	"image1",
	"image2",
	"image3",
	"image4",
	"image5",
}

var smallDiamondHeaders = []string{
	"size_from",
	"size_to",
	"price",
	"quantity",
}

//  "BR" :"圆形" /_images/constant/ico-stones.png - round
//  "PS": "梨形" /_images/constant/ico-stones.png - pear
// 	"PR": "公主方" /_images/constant/ico-stones.png - princess
// 	"HS": "心形" /_images/constant/ico-stones.png - heart
// 	"MQ": "橄榄形" /_images/constant/ico-stones.png - marquise
// 	"OV": "椭圆形" /_images/constant/ico-stones.png - oval
// 	"EM": "祖母绿" /_images/constant/ico-stones.png - emerald
// 	"RA": "雷蒂恩" /_images/constant/ico-stones.png???? - "RAD" - radiant
// 	"CU": "枕形" /_images/constant/ico-stones.png - cushion
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

var VALID_SUPPLIER_NAME = []string{
	"KGK",
	"DIAM",
	"BEYOU-HESHI",
	"SUN",
	"HKEXPO",
	"HESHI",
}

// $clarity_number='0';
/// case "VVS1":
// $clarity_number='1';
/// case "VVS2":
// $clarity_number='2';
// case "VS1":
// $clarity_number='3';
// case "VS2":
// $clarity_number='4';
/// case "SI1":
// $clarity_number='5';
/// case "SI2":
// $clarity_number='6';
/// default:
// $clarity_number='-';
// var VALID_CLARITY_NUMBER = []string{
// 	"0", "1", "2", "3", "4", "5", "6"
// }
// if($dia_clarity=='FL' || $dia_clarity=='IF' || $dia_clarity=='VVS1' || $dia_clarity=='VVS2' || $dia_clarity=='VS1' || $dia_clarity=='VS2'){
// 	$diacomment='瑕疵肉眼不可见';
// }else{
// 	$diacomment='';
// }
var VALID_CLARITY = []string{
	"VVS1",
	"VVS2",
	"VS1",
	"VS2",
	"SI1",
	"SI2",
	"I1",
	"I2",
	"I3",
	"IF",
	"FL",
	"P1",
}

// N - NON -NONE - "None",
// F- FNT - "Faint",
// M - MED - "Medium",
// S - STG - "Strong",
// VST - "Very Strong",
// SL - SLT - "Slight",
// VSL - "Very Slight",
var VALID_FLUORESCENCE_INTENSITY = []string{
	"NON",
	"FNT",
	"MED",
	"STG",
	"VST",
	"SLT",
	"VSL",
	// "NIL",
	// "STR",
	// "VSTG",
}

// "EX" - "EXC" - "Excellent"
// "GD" - "G"
var VALID_POLISH = []string{
	"EX",
	"VG",
	"G",
}

//  "EXC" - "Excellent"-	"EX",
// 	"VG",
// 	"GD" - "G",
// 	"FAIR" - "F",
var VALID_SYMMETRY = []string{
	"EX",
	"VG",
	"G",
	"F",
}

// case "EX":
// $cut_number='0';
// case "VG":
// $cut_number='1';
// case "G":
// $cut_number='2';
// case "F":
// $cut_number='3';
// default:
// $cut_number='-';
// var VALID_CUT_NUMBER = []string{
// 	"0", "1", "2", "3",
// }

//  "EXC" - "Excellent"-	"EX",
// 	"VG",
// 	"GD" - "G",
// 	"FAIR" - "F",
var VALID_CUT_GRADE = []string{
	"EX",
	"VG",
	"G",
	"F",
	// "3EX",
	// "NA",
	// "NN",
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

var VALID_DIAMOND_STATUS = []string{
	"SOLD",
	"AVAILABLE",
	"RESERVED",
	"OFFLINE",
}
