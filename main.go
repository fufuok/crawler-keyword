package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fufuok/crawler-keyword/common"
	"github.com/fufuok/crawler-keyword/conf"
	"github.com/fufuok/crawler-keyword/libs"
	"github.com/fufuok/crawler-keyword/services"
	"github.com/fufuok/crawler-keyword/views"
)

func startColly() {
	for {
		services.Fetch()
		num := libs.RandInt(conf.Config.IntervalMin, conf.Config.IntervalMax)
		common.Info.Println("A 休息", num, "秒")
		time.Sleep(time.Duration(num) * time.Second)
	}
}

func startSelenium() {
	time.Sleep(30 * time.Second)
	for {
		if err := services.FetchGOV(); err != nil {
			common.Error.Println("B 采集器异常:", err)
		}
		num := libs.RandInt(conf.Config.IntervalMin, conf.Config.IntervalMax)
		common.Info.Println("B 休息", num, "秒")
		time.Sleep(time.Duration(num) * time.Second)
	}
}

func envClear() {
	// 创建监听退出 chan
	c := make(chan os.Signal)

	// 监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("程序退出, 清理环境", s)
				command := exec.Command("taskkill", "/T", "/F", "/IM", "ffox.exe")
				err := command.Run()
				if err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Println("程序退出", s)
			}
			os.Exit(0)
		}
	}()
}

func main() {
	envClear()

	if err := conf.InitConfig(""); err != nil {
		common.Error.Panicf("配置加载失败, %s", err)
	}

	// 采集任务
	go startColly()
	go startSelenium()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Static("/data", "./data")

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, views.IndexHmtl())
	})

	_ = r.Run(":21777")
}
