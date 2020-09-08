package goutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"tools/otheroper"
)



func WriteTxt(fileName string,txt string){
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_,err=f.Write([]byte(txt))

	}
}
//复制文件
func Copy(src, dst string) (int64, error) {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	return nBytes, err
}

// 根据 文件 路径 找到文件，并打开当前阅读进度
func GetbookBypath(path string, progress float64) string {
	//打开文件
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)

	defer func() { file.Close() }()

	 ErrorMsg(err, 2, "oper file error:")

	//Stat返回一个描述name指定的文件对象的FileInfo
	infos, _ := os.Stat(path)
	size := infos.Size()
	//当前用户已读文件大小
	currentProgress := progress * float64(size)
	fmt.Printf("当前文件进度为: %v%%, 文件有 %v 个字节 ,文件总大小为 %.2f mb \n", progress*float64(100), size, float64(size)/float64(1048576))

	by := []byte{2}
	//Seek设置下一次读/写的位置。offset为相对偏移量，
	// 而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾。
	// 它返回新的偏移量（相对开头）和可能的错误。
	var bf strings.Builder
	file.Seek(int64(currentProgress), io.SeekStart)
	for i := 0; i <= 1000; i++ {
		file.Read(by)
		bf.Write(by)
	}

	return bf.String()
}

//比较文件类型是否匹配 其实最好的方法是有 文件头来 比较文件类型
func Comperext(filename string, fileext string) string {
	var result string

	if fileext != "*" {
		//如果文件后缀 等于指定文件的类型
		if fileext == filepath.Ext(filename) {
			result = filename
		}
	} else {
		result = filename
	}
	return result
}

//要不要返回目录
func GetContanisDir(needDir bool, f os.FileInfo, fileext string) string {
	var result string

	//要目录
	if needDir {
		//是个目录
		if f.IsDir() {
			result = f.Name()
		} else {
			result = Comperext(f.Name(), fileext)
		}

	} else {
		//不要目录
		if !f.IsDir() {
			//是个文件
			result = Comperext(f.Name(), fileext)

		}
	}
	return result
}

//获取当前目录下的文件，
// path 文件夹路径
// dept 0 当前文件夹下的文件
//      1 当前文件夹所有文件  包含子文件夹下面的文件
// fileext 获取指定文件的类型，如 文本文件 .txt, 任意文件  *
// needDir 要不要返回目录 信息
func GetAllFiles(path string, dept int, fileext string, needDir bool) []string {
	filesSlice := make([]string, 0)

	switch dept {
	case 0:

		infos, _ := ioutil.ReadDir(path)
		for _, f := range infos {

			result := GetContanisDir(needDir, f, fileext)
			if len(result) > 0 {
				filesSlice = append(filesSlice, result)
			}

		}

	case 1:
		root := 0
		//walk 函数会把 root 根路径包含在里面，
		// 就像 tree一样 因此第一次进来时就不返回root了
		err := filepath.Walk(path,
			func(filename string, info os.FileInfo, err error) error {
				otheroper.ErrorMsg(err, 2, "访问文件错误")
				if root > 0 {
					result := GetContanisDir(needDir, info, fileext)
					if len(result) > 0 {
						filesSlice = append(filesSlice, result)
					}
				}
				root = 1

				return nil
			})
		otheroper.ErrorMsg(err, 2, "访问文件错误")

	}

	return filesSlice

}
