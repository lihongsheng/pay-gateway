package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// HttpServer 实现 Server 接口
type HttpServer struct {
	server *http.Server
}

// NewHttpServer 创建 HTTP 服务实例
func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Start 启动 HTTP 服务（实现 Server 接口）
func (h *HttpServer) Start(ctx context.Context) error {
	log.Printf("HTTP 服务启动中，监听地址：%s", h.server.Addr)
	// 使用 goroutine 启动服务，避免阻塞
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("HTTP 服务启动失败：%v", err)
	}
	return nil
}

// Stop 优雅停止 HTTP 服务（实现 Server 接口）
func (h *HttpServer) Stop(ctx context.Context) error {
	log.Println("开始优雅停止 HTTP 服务...")
	// 调用 http.Server 的 Shutdown 方法实现优雅停止
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务停止失败：%w", err)
	}
	log.Println("HTTP 服务已优雅停止")
	return nil
}

// 健康检查接口
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

// 测试接口
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	_, _ = w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
}
