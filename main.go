package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"strconv"
	"strings"
)

type DateInfo struct {
	Date  string
	title string
	Field []string
	Table [][]string
}

func main() {
	var data []*DateInfo
	c := colly.NewCollector()

	baseUrl := "https://www.sge.com.cn"
	listUrl := baseUrl + "/sjzx/mrhqsj?p="

	//定位，并爬取每页上每日的数据
	c.OnHTML("body > div.jzk_main > div > div.jzk_newsCenter_Cont > div.jzk_newsCenter_ContRight"+
		" > div.articleList.border_ea.mt30.mb30 > ul", func(listE *colly.HTMLElement) {

		listE.ForEach("li", func(i int, nodeE *colly.HTMLElement) {
			dayUrl := nodeE.ChildAttr("a", "href") //得到每日行情的链接
			title := nodeE.ChildText("a > span.txt.fl")
			date := nodeE.ChildText("a > span.fr")

			var dateInfo = &DateInfo{Date: date, title: title}
			err := dateInfo.getInfo(baseUrl+dayUrl, c) //通过链接，爬取每日行情数据
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			data = append(data, dateInfo)

		})
	})
	//循环爬取n页的数据
	n:=307
	for i := 1; i <= n; i++ {
		pageListUrl := listUrl + strconv.Itoa(i)
		c.Visit(pageListUrl)
	}
	//后面部分可以选择存进数据库
	bytes,_:=json.Marshal(data)
	ioutil.WriteFile("shanghaiGold.json",bytes,0644)

}

func (dthis *DateInfo) getInfo(url string, srcC *colly.Collector) (err error) {
	c := srcC.Clone()
	dthis.Field = []string{}
	dthis.Table = [][]string{}

	//过滤字符串的\n,\t
	var fs = func(src string) string {
		res:=strings.Replace(src, "\n", "", -1)
		res=strings.Replace(res, "\t", "", -1)
		return res
	}


	//定位到表格数据
	c.OnHTML("body > div.jzk_main > div > div.content.center1200.bgfff > "+
		"div.jzk_newsCenter_meeting.pl30.pr30.pb30 > div.content > table > tbody", func(tableE *colly.HTMLElement) {
			//定位每一行
		tableE.ForEach("tr", func(row int, trE *colly.HTMLElement) {
			if row == 0 {//标题行
				trE.ForEach("td", func(_ int, tdE *colly.HTMLElement) {
					dthis.Field = append(dthis.Field, fs(tdE.Text))
				})
			} else { //数据行
				var tmpArr []string
				trE.ForEach("td", func(_ int, tdE *colly.HTMLElement) {

					tmpArr = append(tmpArr, fs(tdE.Text))

				})
				dthis.Table = append(dthis.Table, tmpArr)

			}
		})

	})
	err = c.Visit(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return nil

}
