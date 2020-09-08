package goutils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"

)

var db *sql.DB

//该dbhelper 没有加入sql注入【验证】
//初始化数据库连接 建议放在main init()中
func InitDB() (err error) {

	dsn := "root:123123@tcp(127.0.0.1:3306)/animationBooks?charset=utf8"

	db, err = sql.Open("mysql", dsn)

	ErrorMsg(err, 2, "open database fail !!")

	err = db.Ping()
	ErrorMsg(err, 2, "ping error:")
	fmt.Println("conn database success")
	return err

}

//反射sql所需的基础模型
func getMode(obj interface{}) ModeObj {
	var modeobj = new(ModeObj)
	*modeobj, _ = GetStructMsg(obj)
	return *modeobj
}

//通用添加方法
func AddStruct(obj interface{}) bool {

	modeobj := getMode(obj)

	sqlStr := fmt.Sprintf("insert into %v values(%v)", modeobj.StructName, modeobj.ValueStrs)
	// stem, err := db.Prepare(sqlStr)
	res, _ := db.Exec(sqlStr)
	//sql 注入
	count, _ := res.RowsAffected()
	return count > 0
}

//批量添加单表数据 主键必须为0
func BatchAddStruct(models []interface{}) bool {
	modeobj := getMode(models)

	var bf strings.Builder

	bf.WriteString(fmt.Sprintf("insert into %v values", modeobj.StructName))
	for _, v := range models {
		modeobj, _ =GetStructMsg(v)
		bf.WriteString(fmt.Sprintf("(%v),", modeobj.ValueStrs))
	}
	sb := bf.String()
	//_, err := db.Prepare(sb[0 : len(sb)-1]+";")
	//fmt.Println(err.Error())
	//开启事务
	//tx, _ := db.Begin()
	//fmt.Println(sb[0 : len(sb)-1],";")

	res, _ := db.Exec(sb[0:len(sb)-1] + ";")
	//fmt.Println(sb[0 : len(sb)-1]+";")
	//sql 注入
	count, _ := res.RowsAffected()
	if count > 0 {
		//	tx.Commit()
	} else {
		//tx.Rollback()
	}
	return count > 0
}

//修改结构
func ModifyStruct(obj interface{}, sqlwhere string, paras ...interface{}) bool {
	modeobj := getMode(obj)

	sqlStr := fmt.Sprintf(" update %v set %v ", modeobj.StructName, sqlwhere)
	res, _ := db.Exec(sqlStr, paras...)

	//sql 注入
	count, _ := res.RowsAffected()

	return count > 0
}

//通用分页方法（单表分页）
func Paging(obj interface{}, cols string, sqlwhere string, primaryKey string, pageIndex int, pageCount int, paras ...interface{}) []map[string]interface{} {
	modeobj := getMode(obj)

	if cols != "*" {
		modeobj.Field = cols

	}

	sqlstr := fmt.Sprintf(" select %v from %v  %v %v >=(select %v from %v limit %v, 1) limit %v;",
		modeobj.Field, modeobj.StructName, sqlwhere, primaryKey, primaryKey, modeobj.StructName, pageIndex*pageCount, pageCount)
	return getlist(sqlstr, paras...)
}

//通用基本查询方法(单表 返回全部字段，或者部分字段)
/*
参数
obj 数据库模型
cols 返回的列名 * 代表全部，
sqlwhere 查询 where条件
paras 多个参数

*/
func QueryStruct(obj interface{}, cols string, sqlwhere string, paras ...interface{}) []map[string]interface{} {

	modeobj := getMode(obj)
	if cols != "*" {
		modeobj.Field = cols

	}
	sqlstr := fmt.Sprintf(" select %v from %v  %v ;", modeobj.Field, modeobj.StructName, sqlwhere)

	return getlist(sqlstr, paras...)
}

//返回对应model的 map 集合
func getlist(sqlstr string, paras ...interface{}) []map[string]interface{} {
	stmt, _ := db.Prepare(sqlstr)
	defer func() { stmt.Close() }()
	/*
	   paras ...interface{}  这里如果添加 可变参数，为一元slice ，在Query 查询时，就变为 二元slice 解析不出来，报错
	*/
	rows, err := stmt.Query(paras...)
	ErrorMsg(err, 2, "")
	defer func() { rows.Close() }()
	var columns, e = rows.Columns()
	ErrorMsg(e, 2, "")
	//https://blog.csdn.net/weimingjue/article/details/91042649 代码出处
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {
		//为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]interface{})
		for i, data := range cache {

			if reflect.TypeOf(*data.(*interface{})).Kind() == reflect.Slice {
				item[columns[i]] = SliceToStringByintefce(*data.(*interface{}))
			} else {
				//类型断言，断定data中type=*interface{}
				item[columns[i]] = *data.(*interface{}) //取实际类型
			}

		}
		fmt.Println(item)
		list = append(list, item)
	}

	return list

}

//根据条件统计总数
func QueryCount(obj interface{}, cols string, sqlwhere string, paras ...interface{}) int {
	modeobj := getMode(obj)
	modeobj.Field = cols
	sqlstr := fmt.Sprintf(" select count(%v) from %v  %v ;", modeobj.Field, modeobj.StructName, sqlwhere)
	stmt, _ := db.Prepare(sqlstr)
	defer func() { stmt.Close() }()
	rows := stmt.QueryRow(paras...)
	count := 0
	rows.Scan(&count)
	return count
}

//自定义sql查询
func QuerySql(sql string, paras ...interface{}) []map[string]interface{} {
	return getlist(sql, paras...)
}

//删除通用方法
func DeleteStruct(obj interface{}, sqlwhere string, paras ...interface{}) bool {

	modeobj := getMode(obj)

	sqlstr := fmt.Sprintf(" delete from %v  %v ;", modeobj.StructName, sqlwhere)
	res, _ := db.Exec(sqlstr, paras...)

	count, err := res.RowsAffected()
	ErrorMsg(err, 2, "")

	return count > 0
}
