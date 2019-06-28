package main

import (
	"fmt"
	"net/http"

	"time"

	"zym.com/utils/alimama/tbk"
	"zym.com/utils/comm_const"
	"zym.com/utils/data_conv/number_lib"
	"zym.com/utils/file"
	"zym.com/utils/http_lib"
)

var tbk = tbk_lib.TaoBaoKey{"25339242", "6322f42345c8ede2d301428ce7620dc7", 64952150385}

const (
	KEY     = "key"
	PAGE_NO = "pageNo"
	TEXT    = "text"
	URL     = "url"
	NUM_IID = "num_iid"
	CAT     = "cat"
)

var buyCount = 0
var defaultCount = 0
var commandCount = 0
var searchCount = 0
var catCount = 0

func main() {

	http.HandleFunc("/search", searchCommodity)
	http.HandleFunc("/getTaoCommand", getTaoCommand)
	http.HandleFunc("/default", getDefaultPage)
	http.HandleFunc("/bigImages", getBigImages)
	http.HandleFunc("/panic_buying2", panicBuying2)
	http.HandleFunc("/cat", getCatCommodify)

	err := http.ListenAndServeTLS(":8080", "/etc/letsencrypt/live/goshop.ink/fullchain.pem",
		"/etc/letsencrypt/live/goshop.ink/privkey.pem", nil)
	if err != nil {
		fmt.Println(err)
	}
}

//商品类目ID,获取商品信息
func getCatCommodify(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, CAT, PAGE_NO); ok {
			var pageNo int32 = 1
			number_lib.StrToInt(params[PAGE_NO], &pageNo)
			body, _ := tbk.DefaultCommodityPage(params[CAT], int(pageNo))
			catCount++
			fmt.Println(time.Now().Format(comm_const.TIME_yyyyMMddHHmmss), "->catCount=", searchCount, "cat=", params[CAT])
			returnWrite(w, body)
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("getCatCommodify->http call is error:", err)
	}
}

//获取秒杀(低价)商品信息
func panicBuying2(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, PAGE_NO); ok {
			var pageNo int32 = 1
			number_lib.StrToInt(params[PAGE_NO], &pageNo)
			body, _ := tbk.PanicBuying2(params[KEY], int(pageNo))
			buyCount++
			fmt.Println(time.Now().Format(comm_const.TIME_yyyyMMddHHmmss), "->buyCount", buyCount)
			returnWrite(w, body)
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("panicBuying->http call is error:", err)
	}
}

//获取商品的图片信息
func getBigImages(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, NUM_IID); ok {
			strImagesUrl := tbk.GetBigImages(params[NUM_IID])
			if strImagesUrl != "" {
				//w.Write([]byte(strImagesUrl))
				returnWrite(w, strImagesUrl)
			} else {
				fmt.Println("get big images fail")
			}
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("getBigImages->http call is error:", err)
	}
}

//获取淘口令
func getTaoCommand(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, TEXT, URL); ok {
			cmd := tbk.CreateTaoCommand(params[TEXT], params[URL])
			commandCount++
			fmt.Println(time.Now().Format(comm_const.TIME_yyyyMMddHHmmss), "->commandCount=", commandCount)
			returnWrite(w, cmd)
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("getTaoCommand->http call is error:", err)
	}
}

//关键字获取商品信息(搜索商品)
func searchCommodity(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, KEY, PAGE_NO); ok {
			var pageNo int32 = 1
			number_lib.StrToInt(params[PAGE_NO], &pageNo)
			body, _ := tbk.SearchCouponCommodity(params[KEY], int(pageNo))
			searchCount++
			fmt.Println(time.Now().Format(comm_const.TIME_yyyyMMddHHmmss), "->searchCount=", searchCount, "key=", params[KEY])
			returnWrite(w, body)
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("searchCommodity->http call is error:", err)
	}
}

//默认页商品信息
func getDefaultPage(w http.ResponseWriter, r *http.Request) {
	params, err := http_lib.GetQueryParam(r)
	if err == nil {
		if n, ok := checkParam(params, PAGE_NO); ok {
			var pageNo int32 = 1
			number_lib.StrToInt(params[PAGE_NO], &pageNo)
			cat := getConfig("default_cat")
			body, _ := tbk.DefaultCommodityPage(cat, int(pageNo))
			defaultCount++
			fmt.Println(time.Now().Format(comm_const.TIME_yyyyMMddHHmmss), "->defatult=", defaultCount)
			returnWrite(w, body)
		} else {
			fmt.Println("lack parameter:", n)
		}
	} else {
		fmt.Println("searchCommodity->http call is error:", err)
	}
}

//检测参数
func checkParam(params map[string]string, paramName ...string) (name string, result bool) {
	result = true
	for i := 0; i < len(paramName); i++ {
		if params[paramName[i]] == "" {
			result = false
			name = paramName[i]
			break
		}
	}
	return
}

//返回信息
func returnWrite(w http.ResponseWriter, data string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := w.Write([]byte(data))
	if err != nil {
		fmt.Println("return data error:", data)
	}
}

func getConfig(field string) (v string) {
	v = file.ReadConfig("CAT", field, "conf/config.txt")
	//fmt.Println("cat:", v)
	return
}
