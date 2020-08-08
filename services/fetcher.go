package services

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"github.com/muesli/cache2go"

	"github.com/fufuok/crawler-keyword/common"
	"github.com/fufuok/crawler-keyword/conf"
)

// 分析网页链接, 得到结果
func ParseA(e *colly.HTMLElement) {
	txt := strings.TrimSpace(e.Attr("title"))
	if txt == "" {
		txt = strings.TrimSpace(e.Text)
	}
	for _, v := range conf.Config.Keywords {
		if strings.Contains(txt, v) {
			common.Info.Println("满足条件:", v, txt)
			// 网站 URL
			sURL, err := conf.CtbCollyUrls.Value(e.Request.URL.String())
			if err != nil {
				common.Error.Println("网址链接获取失败1:", err)
			}
			siteURL := sURL.Data().(string)

			// 符合条件的文章链接, 7天后过期
			href := e.Request.AbsoluteURL(e.Attr("href"))

			// 替换本地链接
			if strings.Contains(href, conf.LocalWeb) {
				href = strings.ReplaceAll(href, conf.LocalWeb, "")
				pURL, err := url.Parse(siteURL)
				if err != nil {
					common.Error.Println("网址链接获取失败2:", err)
				} else {
					okUrl, err := pURL.Parse(href)
					if err != nil {
						common.Error.Println("网址链接获取失败3:", err)
					}
					href = okUrl.String()
				}
			}

			common.Articles.Add(href, 7*86400*time.Second, common.Article{
				Title:   txt,
				URL:     href,
				SiteURL: siteURL,
			})
			break
		}
	}
}

func Fetch() {
	common.Info.Println("A 开始采集")

	c := colly.NewCollector(
		// 开启本机debug
		// colly.Debugger(&debug.LogDebugger{}),
		// 防止页面重复下载
		// colly.CacheDir("./data"),
		// colly.Async(true),
		// colly.MaxDepth(1),
		colly.AllowURLRevisit(),
	)

	q, _ := queue.New(
		2,                                           // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.WithTransport(&http.Transport{
		// Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second, // 超时时间
			KeepAlive: 60 * time.Second, // keepAlive 超时时间
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 空闲连接超时
		TLSHandshakeTimeout:   15 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 1 * time.Second,
	})

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// 对于任何域名, 同时只有两个并发请求在请求该域名, 随机延迟 1 秒
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	c.OnError(func(resp *colly.Response, err error) {
		msg := "等待下轮尝试"
		if resp.StatusCode == http.StatusPreconditionFailed {
			reqURL := resp.Request.URL.String()
			conf.CtbSeleniumUrls.Add(reqURL, 0, reqURL)
			_, _ = conf.CtbCollyUrls.Delete(reqURL)
			msg = "已更换为 B 采集方案"
		}
		common.Error.Printf("A 采集失败, %s, %s, %d, URL: %s", msg, err, resp.StatusCode, resp.Request.URL)
	})

	c.OnHTML("a[href]", ParseA)

	c.OnRequest(func(r *colly.Request) {
		common.Info.Println("A 开始抓取:", r.URL)
	})

	c.OnScraped(func(resp *colly.Response) {
		common.Info.Println("A 抓取完成:", resp.StatusCode, resp.Request.URL)
	})

	conf.CtbCollyUrls.Foreach(func(url interface{}, _ *cache2go.CacheItem) {
		if err := q.AddURL(url.(string)); err != nil {
			common.Trace.Fatalln("队列异常:", err)
		}
	})

	// c.Wait()
	if err := q.Run(c); err != nil {
		common.Trace.Fatalln("程序异常退出:", err)
	}

	common.Info.Println("A 采集全部完成")
}
