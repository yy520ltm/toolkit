package goutils

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

)

func GetEncode(r *bufio.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {

		log.Printf("fetch error :%v", err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

//请求地址 返回内容数组流
func ReptilesDo(url string) ([]byte, error) {

	httprequest, err := http.NewRequest(http.MethodGet, url, nil)

	//httprequest.Close = true
	//防止Go传输本身添加 gzip ，则不会获得 ErrUnexpectedEOF
	httprequest.Header.Add("Accept-Encoding", "utf-8")
	httprequest.Header.Set("Accept", " text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	httprequest.Header.Set("Accept-Charset", "utf-8;q=0.7,*;q=0.3")
	httprequest.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	httprequest.Header.Set("Cache-Control", "max-age=0")
	//模拟浏览器端
	httprequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")

	resp, err := http.DefaultClient.Do(httprequest)
	//resp.Header.Set("Content-Type", "text/plain; charset=utf-8")
	ErrorMsg(err, 2, "请求地址出错")

	defer func() { resp.Body.Close() }()
	fmt.Println("HTTP请求成功,状态码为:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		panic("鸡你太美 statusCode：")

	}

	//bodyRead:=bufio.NewReader(resp.Body)
	//e:=GetEncode(bodyRead)
	//utf8Reader:=transform.NewReader(bodyRead,e.NewDecoder())
	content, err := ioutil.ReadAll(resp.Body)
	//content, err := ioutil.ReadAll(utf8Reader)

	if err != nil {

		if strings.Contains(err.Error(), "unexpected EOF") && len(content) != 0 {
			ErrorMsg(err, 2, " 读这个网站 buffer满了导致写入操作被堵住")
			goto next
		}
	}
	return content, err

next:
	fmt.Println("尝试新方法")

	//如果当前这个网站读取流错误 尝试用这种读取方式

	//fmt.Printf("%s", newContent)
	return nil, err
}
