package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gookit/color"
)

type httpServer struct {
	addr string

	done context.CancelFunc
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

func (hs *httpServer) Run(ctx context.Context, hdlFunc func() http.Handler) (done <-chan struct{}) {
	doneCtx, doneCancle := context.WithCancel(context.Background())
	done = doneCtx.Done()
	hs.done = doneCancle

	go hs.run(ctx, hdlFunc)
	return
}

func (hs *httpServer) run(ctx context.Context, hdlFUnc func() http.Handler) {

	// listen
	ln, err := net.Listen("tcp", hs.addr)
	if err != nil {
		log.Fatalln(cError("Unable to start server: ") + err.Error())
	}
	fmt.Println()
	log.Println("Http server started listening on: ", cDebug("[", hs.addr, "]"))
	log.Println("Swagger ui is serving at: ", cDebug("http://127.0.0.1", hs.addr, "/api/docs/index.html"))
	fmt.Println()

	// httpServer
	srv := &http.Server{
		Addr:    hs.addr,
		Handler: hdlFUnc(),
	}

	// 在 goroutine 中启动 httpServer
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("Unable to start server, error: %v", err.Error())
		}
	}()

	// 接收外部 ctx cancel 通知
	<-ctx.Done()
	log.Println(cDebug("Server Shutdown ..."))
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Err: %+v", err)
	}
	log.Println(cInfo("Server Exited"))
	hs.done()
}

var cDebug = color.Debug.Render
var cInfo = color.Info.Render
var cError = color.FgRed.Render
