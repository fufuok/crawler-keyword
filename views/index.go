package views

import (
	"fmt"
	"strings"

	"github.com/muesli/cache2go"

	"github.com/fufuok/crawler-keyword/common"
	"github.com/fufuok/crawler-keyword/conf"
)

func IndexHmtl() string {
	html := `<!DOCTYPE html>
<html lang="zh-CN" xml:lang="zh-CN">
<head>
    <meta charset="utf-8">
    <title>最新消息</title>
    <style>
        body {
            font: 1.2em/1.6 'Arial Black';
            text-align: center;
            background: aliceblue
        }

        input {
            height: 30px;
            padding: 3px 5px
        }

        a, a:link, a:visited {
            color: #252525;
            text-decoration: none;
			font-weight: 300
        }

        h2 {
            color: #4e97d9;
			margin: 0;
			padding: 0
        }
    </style>
</head>
<body>
    <h1>关键字: {_keywords_}</h1>
    {_body_}
</body>
</html>`
	body := ""
	siteName := ""
	siteURL := ""
	siteArea := ""
	keywords := strings.Join(conf.Config.Keywords, ", ")
	common.Articles.Foreach(func(_ interface{}, item *cache2go.CacheItem) {
		article := item.Data().(common.Article)
		value, err := conf.CtbUrls.Value(article.SiteURL)
		if err != nil {
			common.Error.Println("URL 数据未匹配到站点:", err)
		} else {
			siteInfo := value.Data().(conf.UrlConf)
			siteName = siteInfo.Name
			siteURL = siteInfo.URL
			siteArea = siteInfo.Area
		}
		body = fmt.Sprintf("%s<h2>[<a href='%s' target='_blank'>%s</a>] <a href='%s' target='_blank'>%s</a></h2>",
			body, siteURL, siteName, article.URL, article.Title)
	})
	html = strings.ReplaceAll(html, "{_keywords_}", keywords)
	html = strings.ReplaceAll(html, "{_body_}", body)
	return html
}
