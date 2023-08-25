package http

import (
	"context"
	"log"
	"net/http"
	"time"
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

func (hs *httpServer) Run(ctx context.Context, hdl http.Handler) {
	go hs.run(ctx, hdl)
}

func (hs *httpServer) run(ctx context.Context, hdl http.Handler) {

	// httpServer
	srv := &http.Server{
		Addr:    hs.addr,
		Handler: hdl,
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
