package services

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/muesli/cache2go"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"

	"github.com/fufuok/crawler-keyword/common"
	"github.com/fufuok/crawler-keyword/conf"
	"github.com/fufuok/crawler-keyword/libs"
)

func FetchGOV() error {
	common.Info.Println("B 开始采集")

	port := 21776
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver("./ffox.exe"),
	}

	// selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService("./s381.jar", port, opts...)
	if err != nil {
		return err
	}
	defer func() {
		_ = service.Stop()
	}()

	// 本地 WebDriver
	caps := selenium.Capabilities{"browserName": "firefox"}

	// 禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		// "profile.managed_default_content_settings.images": 2, // chrome
		"permissions.default.image": 2,
	}

	// 浏览器参数
	caps.AddFirefox(firefox.Capabilities{
		Prefs: imagCaps,
		Args: []string{
			// 设置无头模式, 在linux下运行需要设置这个参数, 否则会报错
			"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0",
		},
	})

	// 同步执行, 避免多开浏览器占资源
	conf.CtbSeleniumUrls.Foreach(func(u interface{}, sURL *cache2go.CacheItem) {
		browser, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			return
		}
		// defer wd.Quit()

		urlStr := u.(string)
		filename := libs.MD5(urlStr) + ".html"
		htmlFile := filepath.Join(common.DataDir, filename)
		common.Info.Println("B 开始抓取:", urlStr)

		// 请求网页
		if err := browser.Get(urlStr); err != nil {
			common.Error.Println("B 访问网址出错:", err)
			return
		}

		// 有了先决条件, 再次请求
		if err := browser.Get(urlStr); err != nil {
			common.Error.Println("B 访问网址出错:", err)
			return
		}

		// 获取解析后的网页内容
		body, _ := browser.PageSource()

		// 写入文件
		if err := ioutil.WriteFile(htmlFile, []byte(body), 0644); err != nil {
			common.Error.Println("B 网页内容保存失败:", err)
			return
		}

		// FetchFile(htmlFile)
		// 本地链接加入 Colly, 带原域名地址
		conf.CtbCollyUrls.Add(conf.HtmlWeb+filename, 0, sURL.Data().(string))

		common.Info.Println("B 抓取完成:", htmlFile)

		if err := browser.Quit(); err != nil {
			return
		}
	})

	common.Info.Println("B 采集全部完成")

	return nil
}
