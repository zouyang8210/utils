package tbk_lib

import (
	"fmt"
	"strings"
	"time"
	"utils/comm_const"
	"utils/data_conv/json_lib"
	"utils/http_lib"
)

const (
	SIGN  = "sign"
	EMPTY = ""
	//LAST_PID = 64952150385 //PID最后一组数字

)

//接口名称
const (
	DG_MATERIAL_OPTIONAL = "taobao.tbk.dg.material.optional" //通用物料搜索API（导购）
	TPWD_CREATE          = "taobao.tbk.tpwd.create"
)

//接口URL
const API_GATEWAY = "https://eco.taobao.com/router/rest?"

type TaoBaoKey struct {
	AppKey  string `json:"app_key"`
	Secret  string `json:"secret"`
	LastPid int    `json:"last_pid"` //PID最后一组数字
}

//抢购接口,获取抢购商品
func (sender *TaoBaoKey) PanicBuying(pageNo int) (jsonData string, objData PanicBuyCommodity) {
	params := sender.putPublic("taobao.tbk.ju.tqg.get")
	sender.putParam(params, "adzone_id", sender.LastPid)
	sender.putParam(params, "page_no", pageNo)
	sender.putParam(params, "start_time", time.Now().Format(comm_const.TIME_yyyyMMdd)+" 00:00:00")
	sender.putParam(params, "end_time", time.Now().Format(comm_const.TIME_yyyyMMdd)+" 23:59:59")
	sender.putParam(params, "page_size", 20)
	sender.putParam(params, "fields", "click_url,pic_url,reserve_price,zk_final_price,total_amount,sold_num,title,category_name,start_time,end_time")
	reqStr, _ := sender.putComplete(params)
	resp, err := http_lib.HttpsSubmit(http_lib.GET, reqStr, EMPTY, nil)
	if err != nil {
		fmt.Println("PanicBuying:", err)
	} else {
		jsonData = sender.conv(resp.BufferBody)
		err = json_lib.JsonToObject(jsonData, &objData)
	}
	return
}

//用搜索方式查询价格低的秒杀商品
func (sender *TaoBaoKey) PanicBuying2(key string, pageNo int) (jsonData string, objData CommodityInfo) {
	params := sender.putPublic(DG_MATERIAL_OPTIONAL)
	sender.putParam(params, "adzone_id", sender.LastPid)
	sender.putParam(params, "q", key)
	sender.putParam(params, "page_no", pageNo)
	sender.putParam(params, "has_coupon", true)
	sender.putParam(params, "sort", "total_sales_des") //排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi），价格（price）
	sender.putParam(params, "need_free_shipment", true)
	sender.putParam(params, "page_size", 10)
	sender.putParam(params, "start_price", 1)
	sender.putParam(params, "end_price", 10)
	reqStr, _ := sender.putComplete(params)
	jsonData, objData, err := sender.search(reqStr)
	if err != nil {
		fmt.Println("PanicBuying2:", err)
	}
	return
}

//
func (sender *TaoBaoKey) SearchCouponCommodity(key string, pageNo int) (jsonData string, objData CommodityInfo) {
	params := sender.putPublic(DG_MATERIAL_OPTIONAL)
	sender.putParam(params, "adzone_id", sender.LastPid)
	sender.putParam(params, "q", key)
	sender.putParam(params, "page_no", pageNo)
	sender.putParam(params, "has_coupon", true)
	sender.putParam(params, "sort", "total_sales_des") //排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi），价格（price）
	sender.putParam(params, "need_free_shipment", true)
	sender.putParam(params, "page_size", 10)
	reqStr, _ := sender.putComplete(params)
	jsonData, objData, err := sender.search(reqStr)
	if err != nil {
		fmt.Println("SearchCouponCommodity:", err)
	}
	return
}

func (sender *TaoBaoKey) DefaultCommodityPage(cat string, pageNo int) (jsonData string, objData CommodityInfo) {
	params := sender.putPublic(DG_MATERIAL_OPTIONAL)

	sender.putParam(params, "adzone_id", sender.LastPid)
	sender.putParam(params, "cat", cat)
	sender.putParam(params, "page_no", pageNo)
	sender.putParam(params, "has_coupon", true)
	sender.putParam(params, "sort", "total_sales_des") //排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi），价格（price）
	sender.putParam(params, "need_free_shipment", true)

	reqStr, _ := sender.putComplete(params)
	jsonData, objData, err := sender.search(reqStr)
	if err != nil {
		fmt.Println("DefaultCommodityPage:", err)
	}
	return
}

/*
功能:获取淘口令
参数:
	text:淘口令显示的标题
	url:淘口令url
返回:淘口令字符串
*/
func (sender *TaoBaoKey) CreateTaoCommand(text, url string) (cmd string) {
	params := sender.putPublic(TPWD_CREATE)
	sender.putParam(params, "text", text)
	if strings.Index(url, "https:") < 0 && strings.Index(url, "http:") < 0 {
		sender.putParam(params, "url", "https:"+url)
	} else {
		sender.putParam(params, "url", url)
	}

	reqStr, _ := sender.putComplete(params)
	resp, err := http_lib.HttpsSubmit(http_lib.GET, reqStr, EMPTY, nil)
	if err != nil {
		fmt.Println("CreateTaoCommand:", err)
	} else {
		jsonData := sender.conv(resp.BufferBody)
		var info TaoCommand
		json_lib.JsonToObject(jsonData, &info)
		if info.Code == 0 {
			cmd = info.Data.Model
		} else {
			fmt.Println("淘口令获取失败:", resp.Body)
		}
	}
	return
}

/*
功能:通过商品ID,获取商品祥情的图片地址
参数:
	numIid:商品ID
返回:图片信息
*/
func (sender *TaoBaoKey) GetBigImages(numIid string) (images string) {
	var iid = NumIid{numIid}
	url := "https://h5api.m.taobao.com/h5/mtop.taobao.detail.getdesc/6.0/?data="
	jsonParam, _ := json_lib.ObjectToJson(iid)
	url += jsonParam

	rsp, err := http_lib.HttpsSubmit(http_lib.GET, url, EMPTY, nil)
	if err != nil {
		fmt.Println("GetBigImages:", err)
	} else {
		images = rsp.Body
	}
	return
}
