package comm_const

//时间格式
const (
	MILLISECOND           = ".000"                   //毫秒
	TIME_yyyyMMddHHmmss   = "2006-01-02 15:04:05"    //年月日时分秒
	TIME_MMyyyyddHHmmss   = "01-2006-02 15:04:05"    //月年日时分秒
	TIME_ddMMyyyyHHmmss   = "02-01-2006 15:04:05"    //日月年时分秒
	TIME_yyyyMMddhhmmss12 = "2006-01-02 03:04:05 PM" //-隔符 年-月-日 12小时制
	TIME_MMyyyyddhhmmss12 = "01-2006-02 03:04:05 PM" //-隔符 月-日-年 12小时制
	TIME_ddMMyyyyhhmmss12 = "02-01-2006 03:04:05 PM" //-隔符 日-月-年 12小时制
	TIME_UTC_FORMAT       = "2006-01-02T15:04:05.000000Z07:00"
	TIME_yyyyMMdd         = "2006-01-02"
	TIME_HHmmss           = "15:04:05"
	TIME_hhmmss           = "03:04:05 PM"
)
