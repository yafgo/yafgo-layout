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
	httppkg "yafgo/yafgo-layout/pkg/http"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type WebService struct {
	logger *ylog.Logger
	cfg    *ycfg.Config
}

func NewWebService(logger *ylog.Logger, cfg *ycfg.Config) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,
	}
}

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
	// cfg := g.Cfg()
	cfg := s.cfg
	isProd := cfg.GetString("env") == "prod"
	port := cfg.GetInt("http.port")
	addr := fmt.Sprintf(":%d", port)

	httppkg.NewServerHttp().
		SetAddr(addr).
		Run(ctx, func() http.Handler {
			return NewGinEngine(isProd)
		})
}

var cDebug = color.Debug.Render
var cInfo = color.Info.Render
