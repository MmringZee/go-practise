package apiserver

import (
	"errors"
	mw "fastgo/internal/pkg/middleware"
	genericoptions "fastgo/pkg/options"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	MySQLOptions *genericoptions.MySQLOptions
	Addr         string
}

// Server 定义一个服务器结构体类型.
type Server struct {
	cfg *Config
	srv *http.Server
}

// Run 运行应用.
func (s *Server) Run() error {
	//slog.Info("Read MySQL host from config", "mysql.addr", s.cfg.MySQLOptions.Addr)
	//fmt.Printf("Read MySQL host from config: %s\n", s.cfg.MySQLOptions.Addr)
	//select {} //调用 select 语句，阻塞防止进程退出

	slog.Info("Start to listening the incoming requests on http address", "addr", s.cfg.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// NewServer 根据配置创建服务器.
func (cfg *Config) NewServer() (*Server, error) {
	// 创建gin引擎.
	engine := gin.New()

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.RequestID()}
	// Use()函数入参接收一个可变参数, 但mws是一个切片, 切片后加...可以将其解构为多个独立参数
	engine.Use(mws...)

	// 注册 404 Handler.
	engine.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"code":    "PageNotFound",
			"message": "Page not found",
		})
	})

	// 注册 /healthz handler.
	// 请求方法: GET; 请求路径: /healthz; 请求返回: {"status":"ok"}
	engine.GET("/healthz", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 创建 HTTP Server 实例.
	// 将cfg配置和http服务器实例注入到新创的结构体中
	httpsrv := &http.Server{Addr: cfg.Addr, Handler: engine}

	return &Server{
		cfg: cfg,
		srv: httpsrv,
	}, nil
}
