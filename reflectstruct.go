package goutils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ModeObj struct {
	Field      string //多个字段字符串，以,分隔
	StructName string //结构名称、
	ValueStrs  string //多个值的字符串，以,分隔
}

//根据当前查询出来的元素名称，找到对应结构里面的属性，返回该属性组的 元素切片集合
func GetStructSliceByColumns(value interface{}, colums []string) (valueStructs []interface{}) {
	//创建元素切片
	sliceValues := make([]interface{}, len(colums))
	//取当前的元素类型
	elems := reflect.ValueOf(value)

	for i, v := range colums {

		//参考 https://blog.csdn.net/pkueecser/article/details/50422533
		sliceValues[i] = elems.FieldByName(v).Interface()

	}

	return sliceValues
}

//mp map原始数据
//model 要转换的结构(有问题 )
func ConvertMapToInterface(mp map[string]interface{}, model interface{}) {

	v := reflect.Indirect(reflect.ValueOf(&model))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {

		k := t.Field(i).Name
		var v1 = reflect.ValueOf(v.Field(i))
		v1.Set(reflect.ValueOf(mp[k]))
		//if v.Field(i).Kind()==reflect.Int {
		//	*(*int)(unsafe.Pointer(v.Field(i).Addr().Pointer()))=7
		//}
		//v.FieldByName(k).SetInt(99)
	}

}

//获取结构多个字段字符串,以, 分隔，以及返回结构名称
func GetStructMsg(value interface{}) (modeobj ModeObj, err error) {

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if t.NumField() == 0 {
		return modeobj, errors.New(" you struct Field  length not nil")
	}
	//写入字段
	var bf strings.Builder
	//写入值
	var bfv strings.Builder
	for i := 0; i < t.NumField(); i++ {
		// strings.Join([]string{t.Field(i).Name} ,",")
		bf.WriteString(t.Field(i).Name)
		bf.WriteString(",")
		c := v.Field(i)
		//获取值
		if c.Kind() != reflect.Invalid {

			switch c.Kind() {
			case reflect.String:
				bfv.WriteString(`"`)
				bfv.WriteString(c.String())
				bfv.WriteString(`"`)

			case reflect.Int:
				bfv.WriteString(strconv.FormatInt(c.Int(), 10))

			case reflect.Float32:
				fl32 := strconv.FormatFloat(c.Float(), 'f', -1, 32)
				bfv.WriteString(fl32)

			case reflect.Float64:
				fl64 := strconv.FormatFloat(c.Float(), 'f', -1, 64)

				bfv.WriteString(fl64)

			}

			bfv.WriteString(",")

		}

	}
	bs := bf.String()
	modeobj.Field = bs[0 : len(bs)-1]
	modeobj.StructName = t.Name()

	bs = bfv.String()
	modeobj.ValueStrs = bs[0 : len(bs)-1]
	return modeobj, nil

}
