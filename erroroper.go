package goutils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

//depth o 为 当前ErrorMsg ,2为错误产生的函数
func ErrorMsg(err error, depth int, msg string) bool {

	if err != nil {
		caller(depth, err.Error()+";"+msg)
		return false
	}

	return true
}

//Caller报告当前go程调用栈所执行的函数的文件和行号信息。实参skip为上溯的栈帧数，0表示Caller的调用者
func caller(depth int, msg string) {
	_, file, line, ok := runtime.Caller(depth)

	if !ok {
		msg = "暂未获取到函数信息"
	} else {
		msg = fmt.Sprintf("当前方法在 %v 文件中,在 %v 行中. \n错误信息为: 【%v】  \n", filepath.Base(file), line, msg)
	}

	fmt.Printf(msg)
}
