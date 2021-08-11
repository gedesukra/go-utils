package goutils

import (
	"log"
	"time"
)

func GetDiffDateInSecond(param string) int64 {
	newParam := param + " +0700 WIB"
	format := "2006-01-02 15:04:05 -0700 WIB"
	t, err := time.Parse(format, newParam)
	if err != nil {
		return int64(0)
	}

	delta := time.Now().Sub(t)
	sec := int64(delta.Seconds())

	return sec
}

func GetDiffDateInMinute(param string) int64 {
	newParam := param + " +0700 WIB"
	format := "2006-01-02 15:04:05 -0700 WIB"
	t, err := time.Parse(format, newParam)
	if err != nil {
		return int64(0)
	}

	delta := time.Now().Sub(t)
	min := int64(delta.Minutes())

	return min
}

func GetDiffDateInHour(param string) int64 {
	newParam := param + " +0700 WIB"
	format := "2006-01-02 15:04:05 -0700 WIB"
	t, err := time.Parse(format, newParam)
	if err != nil {
		return int64(0)
	}

	delta := time.Now().Sub(t)
	hour := int64(delta.Hours())

	return hour
}

func GetY() string {
	currentTime := time.Now()
	return currentTime.Format("2006")
}

func GetYMD() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func GetYMDTrans() string {
	currentTime := time.Now()
	return currentTime.Format("20060102")
}

func GetYMDHms() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

func GetYMDHms_ID() string {
	currentTime := time.Now()
	return currentTime.Format("02-01-2006 15:04:05")
}

func GetYMDHmsNoPad() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func GetDDMMMYYYYfromYMDHMS(param string) string {
	layout := "2006-01-02T15:04:05Z"
	t, _ := time.Parse(layout, param)
	return t.Format("02 Jan 2006")
}

func GetDDMMMYYYYHMSfromYMDHMS(param string) string {
	layout := "2006-01-02T15:04:05Z"
	t, _ := time.Parse(layout, param)
	return t.Format("02 Jan 2006 15:04:05")
}

func GetWibfromUTC(param time.Time) time.Time {
	WIB, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
	}
	return param.In(WIB)
}

func GetDiffMinutes(param time.Time) int64 {
	delta := time.Now().Sub(param)
	min := int64(delta.Minutes())
	return min
}

func GetDiffSeconds(param time.Time) int64 {
	delta := time.Now().Sub(param)
	min := int64(delta.Seconds())
	return min
}

func IsGreaterThan(start string, end string, typedate string) bool {
	var layout string
	if typedate == "datetime" {
		layout = "2006-01-02 15:04:05"
	} else {
		layout = "2006-01-02"
	}
	dtStart, _ := time.Parse(layout, start)
	dtEnd, _ := time.Parse(layout, end)
	return dtStart.Before(dtEnd)
}

func IsGreaterThanNow(param string, typedate string) bool {
	var layout string
	if typedate == "datetime" {
		layout = "2006-01-02 15:04:05"
	} else {
		layout = "2006-01-02"
	}
	t, _ := time.Parse(layout, param)
	t1 := time.Now()
	return t1.Before(t)
}

func IsNowBetweenDate(datestart string, dateend string, typedate string) bool {
	var layout string
	if typedate == "datetime" {
		layout = "2006-01-02 15:04:05 -0700 WIB"
	} else {
		layout = "2006-01-02"
	}

	start, _ := time.Parse(layout, datestart+" +0700 WIB")
	end, _ := time.Parse(layout, dateend+" +0700 WIB")
	now := time.Now()
	validGreaterEquals := now.After(start) || now.Equal(start)
	validLowerEquals := now.Before(end) || now.Equal(end)

	return validGreaterEquals && validLowerEquals
}
