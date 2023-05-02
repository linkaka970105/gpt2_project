package util

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 定时一些常量
const (
	DaySeconds int64 = 3600 * 24
)

// ZeroTsWithLoc 今天零点的时间戳
func ZeroTsWithLoc(t time.Time, loc *time.Location) int64 {
	t = t.In(loc)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Unix()
}

// ZeroTs 零点的时间戳
func ZeroTs(t time.Time) int64 {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Unix()
}

// ZeroWithLoc 今天零点
func ZeroWithLoc(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// JSON json encode
func JSON(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

// UUID new uuid
func UUID() uuid.UUID {
	return uuid.New()
}

// JoinIntSlice 把[]int转换为字符串
func JoinIntSlice(a []int, sep string) string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

// IsSameDay 判断是否同一天
func IsSameDay(t1 time.Time, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// DateStr 时间转日期字符串
func DateStr(t time.Time) string {
	return t.Format("2006-01-02")
}

// DateTimeStr 时间转日期时间字符串
func DateTimeStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

const TimeFormat = "2006-01-02 15:04:05"
const TimeMinuteFormat = "2006-01-02 15:04"
const TimeDateFormat = "2006-01-02"

// 将<br>标签替换为\n，并去掉其他h5标签
func H5TagReplace(text string) string {
	if text == "" {
		return text
	}
	plainText := strings.ReplaceAll(text, "<br>", "\n")
	regex, _ := regexp.Compile("<[^>]+>")
	plainText = regex.ReplaceAllString(plainText, "")
	return plainText
}

func BinaryToIntArray(n int) []int {
	var result []int
	for i := 0; n > 0; i++ {
		if n%2 == 1 {
			result = append(result, 1<<i)
		}
		n = n >> 1
	}
	return result
}

func IntArrayToBinary(arr []int) int {
	var result int
	for _, n := range arr {
		result |= n
	}
	return result
}

func log2(n int) int {
	var count int
	for n > 1 {
		n >>= 1
		count++
	}
	return count
}

// 1,2,4,8 ==> 1,2,3,4
func Int2Index(input []int) (out []int) {
	for _, n := range input {
		out = append(out, log2(n)+1)
	}
	return
}

// 1,2,3,4 ==> 1,2,4,8
func Index2Int(input []int) (out []int) {
	for _, n := range input {
		out = append(out, 1<<(n-1))
	}
	return
}
