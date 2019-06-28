// time_lib
package time_lib

import (
	"fmt"
	"time"
)

//时间格式24小时制
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

//本地时间
const LOCAL = "Local"

func TimeToStr(t time.Time, format string) string {
	return t.Format(format)
}

/*
功能：时间字符串转换成time.Time（本地时间）
参数：
	s_time:时间字符串
返回：转换的时间
*/
func StrToTime(strTime string) (t time.Time) {
	loc, _ := time.LoadLocation(LOCAL)
	t, _ = time.ParseInLocation(TIME_yyyyMMddHHmmss, strTime, loc)
	//t, _ = time.Parse(TIME_yyyyMMddHHmmss, strTime)
	return
}

/*
功能：时间加上或减去 N 小时、分钟或秒
参数：
	dt:需要计算的时间
	num:数字 负数表示减
	falg:需要计算的单位（h=小时 m=分钟 s=秒
返回：计算后的时间
*/
func TimeAdd(dt time.Time, num int64, flag string) (t time.Time) {
	m, _ := time.ParseDuration(fmt.Sprintf("%d%s", num, flag))
	t = dt.Add(m)
	return
}

/*
功能：计算两个时间相差的秒数
参数：
	t1:时间1
	t2:时间2
返回：两时间相差的秒数
*/
func TimeSub(t1 time.Time, t2 time.Time) (sec int64) {
	if t1.Before(t2) {
		sec = t2.Unix() - t1.Unix()
	} else {
		sec = t1.Unix() - t2.Unix()
	}
	return
}

/*
功能:把yyyyMMddHHmmss格式的时间,转在标准字符串时间
参数:
	date:yyyyMMddHHmmss格式的字符串
	separator:时间分隔符 1='-',2='/'
返回:
	指定格式日期
*/
func StrDateFormat(date, separator string) (strTime string) {
	year := date[0:4]
	month := date[4:6]
	day := date[6:8]
	hour := date[8:10]
	minute := date[10:12]
	second := date[12:14]

	strTime = fmt.Sprintf("%s%s%s%s%s %s:%s:%s", year, separator, month, separator, day,
		hour, minute, second)
	return
}

/*
功能:字符串时间计算时间截
参数:
	t:字符串时间
返回:时间戳
*/
func StrToTimestamp(t string) int64 {
	loc, _ := time.LoadLocation(LOCAL)
	tm, err := time.ParseInLocation(TIME_yyyyMMddHHmmss, t, loc)
	if err != nil {
		return -1
	}
	return int64(tm.Unix())
}

/*
功能:配合defer,查看代码运行时间
参数:
	now:时间
返回:时间戳
*/
func RunTime(now time.Time, flag string) {
	terminal := time.Since(now)
	fmt.Println(flag, terminal)
}
