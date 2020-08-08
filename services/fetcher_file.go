package services

import (
	"net/http"

	"github.com/gocolly/colly"

	"github.com/fufuok/crawler-keyword/common"
	"github.com/fufuok/crawler-keyword/libs"
)

// 访问本地网页文件
func FetchFile(file string) {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	c.OnHTML("a[href]", ParseA)

	c.OnRequest(func(r *colly.Request) {
		common.Info.Println("B 开始抓取:", r.URL)
	})

	c.OnError(func(resp *colly.Response, err error) {
		common.Error.Printf("B 采集失败, %s, Url: %s", err, resp.Request.URL)
	})

	if err := c.Visit(libs.FileUrl(file)); err != nil {
		common.Error.Println("B 网页文件访问失败:", err)
	}

	c.Wait()
}
