//golang实现的简单爬虫，爬豆瓣电影TOP250的电影名、评分、评价人数等信息，创建写入到制定文件夹下的excel文件中。

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

//定义新的数据类型
type Spider struct {
	url    string
	header map[string]string
}

//定义 Spider get的方法
func (keyword Spider) get_html_header() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", keyword.url, nil)
	if err != nil {
		//handle error
		fmt.Println("error in #1")
	}
	for key, value := range keyword.header {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in #2")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in #3")
	}
	return string(body)

}
func parse() {
	header := map[string]string{
		"Host":                      "movie.douban.com",
		"Connection":                "keep-alive",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1", //让浏览器不再显示https页面中对于http请求的报警
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Referer":                   "https://movie.douban.com/top250",
	}

	//创建excel文件
	f, err := os.Create("C:/Users/ctao8/Desktop/crawler.xls")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//写入标题
	f.WriteString("电影名称" + "\t" + "评分" + "\t" + "评价人数" + "\t" + "\r\n")

	//循环每页解析并把结果写入excel
	for i := 0; i < 10; i++ {
		fmt.Println("正在抓取第" + strconv.Itoa(i) + "页......") //strconv.Itoa:将一个整数转为字符串
    
		/*
      URI规律：
			第一页uri：https://movie.douban.com/top250?start=0&filter=
			第二页uri：https://movie.douban.com/top250?start=25&filter=
		*/
    
		url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter=" //因为uri有规律所以可以循环区抓取
		spider := &Spider{url, header}
		html := spider.get_html_header()

		//评价人数
		//html中关于评价的行：<span>592195人评价</span>
		pattern2 := `<span>(.*?)评价</span>`
		rp2 := regexp.MustCompile(pattern2)            //MustCompile解析并返回一个正则表达式。如果成功返回，该Regexp就可用于匹配文本，主要用于全局正则表达式变量的安全初始化
		find_txt2 := rp2.FindAllStringSubmatch(html, -1)
    
    /*
    test:推断类型
    fmt.Println("type:", reflect.TypeOf(pattern2)) 
    fmt.Println("type:", reflect.TypeOf(rp2))
  */
  
		//评分
		//html中关于评分的字段：property="v:average">9.4</span>
		pattern3 := `property="v:average">(.*?)</span>`
		rp3 := regexp.MustCompile(pattern3)
		find_txt3 := rp3.FindAllStringSubmatch(html, -1)
		//FindAllStringSubmatch:返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的[][]string切片。如果没有匹配到，会返回nil。

		//电影名称
		//html中关于评分的字段：alt="阿甘正传" src=
		pattern4 := `alt="(.*?)" src="`
		rp4 := regexp.MustCompile(pattern4)
		find_txt4 := rp4.FindAllStringSubmatch(html, -1)

/*
test:查看是否有爬到数据
		fmt.Println(len(find_txt2))
		fmt.Println(len(find_txt3))
		fmt.Println(len(find_txt4))
*/


		// 写入UTF-8 BOM
		f.WriteString("\xEF\xBB\xBF")
		//  打印全部数据和写入excel文件
		for i := 0; i < len(find_txt2); i++ {
			//fmt.Printf("%s %s %s\n", find_txt4[i][0], find_txt3[i][0], find_txt2[i][0])
			f.WriteString(find_txt4[i][1] + "\t" + find_txt3[i][1] + "\t" + find_txt2[i][1] + "\t" + "\r\n")
		}

	}
}

func main() {
	t1 := time.Now() // get current time
	parse()
	elapsed := time.Since(t1)
	fmt.Println("爬虫结束,总共耗时: ", elapsed)

}
