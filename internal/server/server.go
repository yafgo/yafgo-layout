package server

import (
	"context"
	"fmt"
	"go-toy/toy-layout/internal/global"
	"go-toy/toy-layout/pkg/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type webService struct{}

func NewWebService() *webService {
	return &webService{}
}

func (s *webService) CmdRun(cmd *cobra.Command, args []string) {

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
		log.Printf("%s, waiting...\n", sig)
		cancel()
	}()

	s.RunWebServer(ctx)

	// 等待退出
	<-ctx.Done()
	// 缓冲几秒等待任务结束
	log.Println("exiting...")
	time.Sleep(time.Second * 2)
	log.Println("exit")
}

// RunWebServer 启动 web server
func (s *webService) RunWebServer(ctx context.Context) {
	cfg := global.Ycfg
	isProd := cfg.GetString("env") == "prod"
	port := cfg.GetInt("http.port")
	addr := fmt.Sprintf(":%d", port)

	ginR := NewGinEngine(isProd)
	http.NewServerHttp().
		SetAddr(addr).
		Run(ctx, ginR)
}
