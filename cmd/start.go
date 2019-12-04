package main

import (
	"context"
	"gin/internal/databases/mysql"
	"gin/internal/util"
	"gin/routers"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	g       errgroup.Group
	servers map[string]*http.Server
)

// home
func routerHome(router *gin.Engine) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	routers.StuffHttp(router)

	// 404页面
	router.NoRoute(htmlPage404)

	return e
}

// api
func routerApi(router *gin.Engine) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	routers.StuffApi(router)

	// 404页面
	router.NoRoute(ApiPage404)

	return e
}

func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	rootPath := os.Getenv("PWD")

	router.LoadHTMLGlob(util.GetParentDirectory(rootPath) + "/app/view/**/*")

	// 加载 home 路由
	home := routerHome(router)

	// 加载 api 路由
	api := routerApi(router)

	defer func() {
		mysql.Close()
	}()

	port := "8888"

	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}

	servers["home"] = &http.Server{
		Addr:         ":" + port,
		Handler:      home,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	servers["api"] = &http.Server{
		Addr:         ":" + port,
		Handler:      api,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	for _, server := range servers {
		//go start(server, router)
		g.Go(func() error {
			return server.ListenAndServe()
		})
	}

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

func start(srv *http.Server, router *gin.Engine) {
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// 404 页面
func htmlPage404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"message": "服务器开小差了，很快就好，稍后再刷新试试吧~",
	})
}

// 404 页面
func ApiPage404(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "服务器开小差了，很快就好，稍后再刷新试试吧~",
	})
}
