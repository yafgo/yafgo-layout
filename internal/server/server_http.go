package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yafgo/yafgo-layout/internal/middleware"
	httppkg "yafgo/yafgo-layout/pkg/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var cDebug = color.Debug.Render
var cInfo = color.Info.Render

func (s *WebService) CmdRun(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())

	// 监听关停信号
	sigs := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 监听外部终止程序的信号
	go func() {
		sig := <-sigs
		log.Printf("%s, waiting...", sig)
		cancel()
	}()

	s.RunWebServer(ctx)

	// 等待退出
	<-ctx.Done()
	// 缓冲几秒等待任务结束
	log.Println(cDebug("exiting..."))
	time.Sleep(time.Second * 2)
	log.Println(cInfo("exit"))
}

// RunWebServer 启动 web server
func (s *WebService) RunWebServer(ctx context.Context) {
	isProd := s.g.IsProd()
	port := s.cfg.GetInt("http.port")
	addr := fmt.Sprintf(":%d", port)

	httppkg.NewServerHttp().
		SetAddr(addr).
		Run(ctx, func() http.Handler {
			return s.NewGinEngine(isProd)
		})
}

// NewGinEngine
func (s *WebService) NewGinEngine(isProd bool) *gin.Engine {
	// 设置 gin 的运行模式，支持 debug, release, test, 生产环境请使用 release 模式
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化 Gin 实例
	r := gin.New()

	// 注册全局中间件
	s.registerGlobalMiddleware(r)

	// 注册路由
	s.registerRoutes(r)

	return r
}

// registerGlobalMiddleware 注册全局中间件
func (s *WebService) registerGlobalMiddleware(router *gin.Engine) {
	router.Use(
		middleware.Logger(),
		gin.Recovery(),
		middleware.CORS(),
		/* func(ctx *gin.Context) {
			ctx.Header("Access-Control-Expose-Headers", "custom-header")
		}, */
	)
}
