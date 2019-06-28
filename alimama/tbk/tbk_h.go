package tbk_lib

type RetBase struct {
	RequestId  string `json:"request_id"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	SubErrCode string `json:"sub_err_code"`
	SubErrMsg  string `json:"sub_err_msg"`
}

type PanicBuyCommodity struct {
	RetBase
	Results PanicBuyList `json:"results"`
}

type PanicBuyList struct {
	Results []PanicBuyItem `json:"results"`
}

type PanicBuyItem struct {
	CategoryName string `json:"category_name"` //商品类目名称
	ClickUrl     string `json:"click_url"`
	EndTime      string `json:"end_time"`
	StartTime    string `json:"start_time"`
	NumIid       int    `json:"num_iid"`
	PicUrl       string `json:"pic_url"`
	ReservePrice string `json:"reserve_price"` //原价
	SoldNum      int    `json:"sold_num"`      //已抢购数量
	Title        string `json:"title"`
	TotalAmount  int    `json:"total_amount"`
	ZkFinalPrice string `json:"zk_final_price"` //折扣价
}

type NumIid struct {
	Id string `json:"id"`
}

//===================淘口令=======================
type TaoCommand struct {
	RetBase
	Data TaoCommandData `json:"data"`
}

type TaoCommandData struct {
	Model string `json:"model"`
}

//============================================================================================================

//==============================搜索商品信息=====================================
type CommodityInfo struct {
	ResultList   CommodityMap `json:"result_list"`
	TotalResults int          `json:"total_results"`
	RetBase
}

type CommodityMap struct {
	MapData []Commodity `json:"map_data"`
}

//设备信息明细
type Commodity struct {
	NumIid         int        `json:"num_iid"`          ////商品ID
	CouponId       string     `json:"coupon_id"`        //优惠券ID
	PictUrl        string     `json:"pict_url"`         //商品主图
	SmallImages    SmallImage `json:"small_images"`     //商品小图列表
	UserType       int        `json:"user_type"`        //卖家类型，0表示集市，1表示商城
	ItemUrl        string     `json:"item_url"`         //商品地址
	CommissionRate string     `json:"commission_rate"`  //佣金比率
	Volume         int        `json:"volume"`           //30天售量
	CouponShareUrl string     `json:"coupon_share_url"` //券二合一页面链接(领优惠券地址)
	Url            string     `json:"url"`              //商品淘宝客链接
	ShortTitle     string     `json:"short_title"`      //商品短标题
	ReservePrice   string     `json:"reserve_price"`    //商品一口价
	ZkFinalPrice   string     `json:"zk_final_price"`   //商品折扣价
	CouponInfo     string     `json:"coupon_info"`      //优惠券面额(满299元减20元)
	CouponLater    float32    `json:"couponLater"`      //领券之后的价格
	CouponPrice    float32    `json:"couponPrice"`      //优惠券价格
	TaoCommand     string     `json:"taoCommand"`       //淘品令
	Title          string     `json:"title"`
}

//小图片数组
type SmallImage struct {
	SmallImgs []string `json:"small_imgs"`
}

//==========================================================================================================
