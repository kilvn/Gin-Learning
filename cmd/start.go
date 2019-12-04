package main

import (
	"gin/internal/databases/mysql"
	"gin/internal/util"
	"gin/routers"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	g       errgroup.Group
	servers map[string]*http.Server
	msg404  = "服务器开小差了，很快就好，稍后再刷新试试吧~"
)

// home
func routerHome(router *gin.Engine) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	// 404页面
	router.NoRoute(htmlPage404)

	routers.StuffHttp(router)

	return e
}

// api
func routerApi(router *gin.Engine) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	// 404页面
	router.NoRoute(ApiPage404)

	routers.StuffApi(router)

	return e
}

func main() {
	gin.SetMode(gin.DebugMode)

	defer func() {
		mysql.Close()
	}()

	router := gin.Default()

	rootPath := os.Getenv("PWD")

	router.LoadHTMLGlob(util.GetParentDirectory(rootPath) + "/app/view/**/*")

	// 加载 home 路由
	homeHandler := routerHome(router)

	// 加载 api 路由
	apiHandler := routerApi(router)

	port := ":8888"

	err := router.Run(port)
	if err != nil {
		panic(err)
	}

	servers["home"] = &http.Server{
		Addr:           port,
		Handler:        homeHandler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	servers["api"] = &http.Server{
		Addr:           port,
		Handler:        apiHandler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	for _, server := range servers {
		g.Go(func() error {
			return server.ListenAndServe()
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// 404 页面
func htmlPage404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"message": msg404,
	})
}

// 404 页面
func ApiPage404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": msg404,
	})
}
