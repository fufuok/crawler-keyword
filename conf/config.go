package conf

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"github.com/muesli/cache2go"

	"github.com/fufuok/crawler-keyword/common"
)

// 网址配置
type UrlConf struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Area string `json:"area"`
}

type JsonConf struct {
	IntervalMin int       `json:"interval_min"`
	IntervalMax int       `json:"interval_max"`
	Keywords    []string  `json:"keywords"`
	URLs        []UrlConf `json:"urls"`
}

// 所有配置
var Config *JsonConf

// 网址为键的配置数据 UrlConf.Url.Parse: UrlConf
var CtbUrls = cache2go.Cache("URLs")

// 待处理网址集 colly: []string, selenium []string
var CtbCollyUrls = cache2go.Cache("CollyUrls")
var CtbSeleniumUrls = cache2go.Cache("SeleniumUrls")

// 本地 Web 服务
var LocalWeb = "http://127.0.0.1:21777"
var HtmlWeb = LocalWeb + "/data/"

// 加载配置
func InitConfig(filename string) error {
	if filename == "" {
		filename = filepath.Join(common.RootDir, "conf.json")
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// demo
	_ = []byte(`{
  "keywords": [
    "GTA5",
    "会议",
    "节点",
    "公告",
    "模式",
    "董事长",
    "通知"
  ],
  "urls": [
    {
      "name": "网站一",
      "url": "https://www.xunyou.com/"
    },
    {
      "name": "网站二",
      "url": "http://www.qbj.gov.cn/"
    },
    {
      "name": "网站三",
      "url": "https://cs.xunyou.com/html/282/"
    }
  ]
}`)

	if err := json.Unmarshal(body, &Config); err != nil {
		return err
	}

	for _, v := range Config.URLs {
		_, err := url.Parse(v.URL)
		if err != nil {
			common.Error.Printf("网址配置有误, [忽略], %s", v.URL)
			continue
		}
		// url.Url 为键的数据, 页面展示时匹配显示
		CtbUrls.Add(v.URL, 0, v)
		CtbCollyUrls.Add(v.URL, 0, v.URL)
	}

	if Config.IntervalMin < 30 {
		Config.IntervalMin = 30
	}
	if Config.IntervalMax < Config.IntervalMin {
		Config.IntervalMax = 60
	}

	return nil
}
