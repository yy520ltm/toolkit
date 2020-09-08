package toolkit

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
	"unsafe"
)

//byte数组转 int32
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

//byte数组转 int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// 实现int64转换成byte数组
func Int64Byte(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	ErrorMsg(err, 2, "int64 convert error!")
	return buffer.Bytes()
}

// 第二种实现int64转换成byte数组
func Int64Byte2(num int64) []byte {
	return []byte(strconv.FormatInt(num, 10))
}


//string 转换成byte[]
func StringTobytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func SliceToString(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

//利用反射来转切片为string
func SliceToStringByintefce(v interface{}) (s string) {
	/*
		if reflect.TypeOf(v).Kind()!=reflect.Slice {
			otheroper.ErrorMsg(errors.New("当前类型不是Slice"),2,"")
			return
		}
	*/

	return string(reflect.ValueOf(v).Bytes())
}


//  字符首字母大写
func FirstUpper(str string) string {
	var upperStr string

	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
