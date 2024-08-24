package constant

import "github.com/hypay-id/backend-dashboard-hypay/internal/dto"

// TransformJackBank is bank code translator from internal bank code into jack bank code
var TransformJackBank = map[string]dto.JackBankCode{
	"IDR_061": {
		Id:       "16",
		BankName: "anz",
	},
	"IDR_116": {
		Id:       "8",
		BankName: "aceh",
	},
	"IDR_088": {
		Id:       "15",
		BankName: "antar_daerah",
	},
	"IDR_037": {
		Id:       "18",
		BankName: "artha",
	},
	"IDR_542": {
		Id:       "153",
		BankName: "jago",
	},
	"IDR_133": {
		Id:       "25",
		BankName: "bengkulu",
	},
	"IDR_547": {
		Id:       "37",
		BankName: "btpn_syar",
	},
	"IDR_521": {
		Id:       "39",
		BankName: "bukopin_syar",
	},
	"IDR_076": {
		Id:       "40",
		BankName: "bumi_artha",
	},
	"IDR_054": {
		Id:       "42",
		BankName: "capital",
	},
	"IDR_014": {
		Id:       "3",
		BankName: "bca",
	},
	"IDR_036": {
		Id:       "46",
		BankName: "china_cons",
	},
	"IDR_011": {
		Id:       "7",
		BankName: "danamon",
	},
	"IDR_111": {
		Id:       "58",
		BankName: "dki",
	},
	"IDR_161": {
		Id:       "62",
		BankName: "ganesha",
	},
	"IDR_567": {
		Id:       "64",
		BankName: "harda",
	},
	"IDR_513": {
		Id:       "68",
		BankName: "ina_perdana",
	},
	"IDR_555": {
		Id:       "69",
		BankName: "index_selindo",
	},
	"IDR_115": {
		Id:       "71",
		BankName: "jambi",
	},
	"IDR_472": {
		Id:       "72",
		BankName: "jasa_jakarta",
	},
	"IDR_113": {
		Id:       "73",
		BankName: "jateng",
	},
	"IDR_114": {
		Id:       "75",
		BankName: "jatim",
	},
	"IDR_123": {
		Id:       "78",
		BankName: "kalbar",
	},
	"IDR_122": {
		Id:       "80",
		BankName: "kalsel",
	},
	"IDR_125": {
		Id:       "82",
		BankName: "kalteng",
	},
	"IDR_124": {
		Id:       "84",
		BankName: "kaltim",
	},
	"IDR_535": {
		Id:       "229",
		BankName: "seabank",
	},
	"IDR_121": {
		Id:       "86",
		BankName: "lampung",
	},
	"IDR_131": {
		Id:       "87",
		BankName: "maluku",
	},
	"IDR_008": {
		Id:       "2",
		BankName: "mandiri",
	},
	"IDR_564": {
		Id:       "162",
		BankName: "mantap",
	},
	"IDR_157": {
		Id:       "90",
		BankName: "maspion",
	},
	"IDR_097": {
		Id:       "91",
		BankName: "mayapada",
	},
	"IDR_553": {
		Id:       "95",
		BankName: "mayora",
	},
	"IDR_426": {
		Id:       "97",
		BankName: "mega_tbk",
	},
	"IDR_506": {
		Id:       "96",
		BankName: "mega_syar",
	},
	"IDR_151": {
		Id:       "98",
		BankName: "mestika",
	},
	"IDR_485": {
		Id:       "103",
		BankName: "mnc",
	},
	"IDR_548": {
		Id:       "105",
		BankName: "multiarta",
	},
	"IDR_095": {
		Id:       "77",
		BankName: "jtrust",
	},
	"IDR_128": {
		Id:       "110",
		BankName: "ntb",
	},
	"IDR_130": {
		Id:       "111",
		BankName: "ntt",
	},
	"IDR_145": {
		Id:       "7",
		BankName: "danamon",
	},
	"IDR_069": {
		Id:       "45",
		BankName: "china",
	},
	"IDR_146": {
		Id:       "70",
		BankName: "india",
	},
	"IDR_042": {
		Id:       "101",
		BankName: "mitsubishi",
	},
	"IDR_132": {
		Id:       "116",
		BankName: "papua",
	},
	"IDR_013": {
		Id:       "6",
		BankName: "permata",
	},
	"IDR_520": {
		Id:       "118",
		BankName: "prima_master",
	},
	"IDR_002": {
		Id:       "4",
		BankName: "bri",
	},
	"IDR_119": {
		Id:       "125",
		BankName: "riau",
	},
	"IDR_523": {
		Id:       "127",
		BankName: "sampoerna",
	},
	"IDR_152": {
		Id:       "129",
		BankName: "shinhan",
	},
	"IDR_153": {
		Id:       "130",
		BankName: "sinarmas",
	},
	"IDR_126": {
		Id:       "133",
		BankName: "sulselbar",
	},
	"IDR_134": {
		Id:       "135",
		BankName: "sulteng",
	},
	"IDR_135": {
		Id:       "136",
		BankName: "sultenggara",
	},
	"IDR_127": {
		Id:       "137",
		BankName: "sulut",
	},
	"IDR_118": {
		Id:       "138",
		BankName: "sumbar",
	},
	"IDR_120": {
		Id:       "140",
		BankName: "sumsel_babel",
	},
	"IDR_117": {
		Id:       "142",
		BankName: "sumut",
	},
	"IDR_451": {
		Id:       "154",
		BankName: "bsi",
	},
	"IDR_566": {
		Id:       "145",
		BankName: "victoria",
	},
	"IDR_405": {
		Id:       "146",
		BankName: "victoria_syar",
	},
	"IDR_068": {
		Id:       "147",
		BankName: "woori",
	},
	"IDR_490": {
		Id:       "148",
		BankName: "yudha_bhakti",
	},
	"IDR_536": {
		Id:       "24",
		BankName: "bca_syar",
	},
	"IDR_110": {
		Id:       "27",
		BankName: "bjb",
	},
	"IDR_425": {
		Id:       "28",
		BankName: "bjb_syar",
	},
	"IDR_009": {
		Id:       "1",
		BankName: "bni",
	},
	"IDR_129": {
		Id:       "20",
		BankName: "bali",
	},
	"IDR_137": {
		Id:       "22",
		BankName: "banten",
	},
	"IDR_112": {
		Id:       "56",
		BankName: "diy",
	},
	"IDR_494": {
		Id:       "100",
		BankName: "mitraniaga",
	},
	"IDR_200": {
		Id:       "34",
		BankName: "btn",
	},
	"IDR_213": {
		Id:       "36",
		BankName: "btpn",
	},
	"IDR_441": {
		Id:       "38",
		BankName: "bukopin",
	},
	"IDR_022": {
		Id:       "5",
		BankName: "cimb",
	},
	"IDR_031": {
		Id:       "50",
		BankName: "citibank",
	},
	"IDR_950": {
		Id:       "51",
		BankName: "commonwealth",
	},
	"IDR_949": {
		Id:       "47",
		BankName: "chinatrust",
	},
	"IDR_046": {
		Id:       "53",
		BankName: "dbs",
	},
	"IDR_041": {
		Id:       "66",
		BankName: "hsbc",
	},
	"IDR_164": {
		Id:       "67",
		BankName: "icbc",
	},
	"IDR_484": {
		Id:       "63",
		BankName: "hana",
	},
	"IDR_016": {
		Id:       "92",
		BankName: "maybank",
	},
	"IDR_147": {
		Id:       "104",
		BankName: "muamalat",
	},
	"IDR_503": {
		Id:       "109",
		BankName: "nobu",
	},
	"IDR_028": {
		Id:       "112",
		BankName: "ocbc",
	},
	"IDR_019": {
		Id:       "114",
		BankName: "panin",
	},
	"IDR_517": {
		Id:       "115",
		BankName: "panin_syar",
	},
	"IDR_167": {
		Id:       "121",
		BankName: "qnb",
	},
	"IDR_498": {
		Id:       "128",
		BankName: "sbi",
	},
	"IDR_050": {
		Id:       "132",
		BankName: "stanchard",
	},
	"IDR_023": {
		Id:       "144",
		BankName: "uob",
	},
	"IDR_945": {
		Id:       "164",
		BankName: "ibk",
	},
	"IDR_201": {
		Id:       "35",
		BankName: "btn_syar",
	},
	"ID_OVO": {
		Id:       "150",
		BankName: "ovo",
	},
	"ID_DANA": {
		Id:       "166",
		BankName: "dana",
	},
	"ID_GOPAY": {
		Id:       "173",
		BankName: "gopay",
	},
	"ID_SHOPEEPAY": {
		Id:       "236",
		BankName: "shopeepay",
	},
	"ID_LINKAJA": {
		Id:       "165",
		BankName: "linkaja",
	},
}
