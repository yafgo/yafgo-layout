package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	addr string
}

func NewServerHttp() *httpServer {
	return &httpServer{
		addr: ":8080",
	}
}

// SetAddr 设置监听端口
func (hs *httpServer) SetAddr(addr string) *httpServer {
	hs.addr = addr
	return hs
}

func (hs *httpServer) Run(ctx context.Context, r *gin.Engine) {
	go hs.run(ctx, r)
}

func (hs *httpServer) run(ctx context.Context, r *gin.Engine) {

	// httpServer
	srv := &http.Server{
		Addr:    hs.addr,
		Handler: r,
	}

	// 在 goroutine 中启动 httpServer
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unable to start server, error: %v", err.Error())
		}
	}()

	// 接收外部 ctx cancel 通知
	<-ctx.Done()
	log.Println("Server Shutdown ...")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Err:", err)
	}
	log.Println("Server Exited")
}
