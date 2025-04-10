package apiserver

import (
	"context"
	"errors"
	"fastgo/internal/apiserver/biz"
	"fastgo/internal/apiserver/handler"
	"fastgo/internal/apiserver/pkg/validation"
	store2 "fastgo/internal/apiserver/store"
	"fastgo/internal/pkg/core"
	"fastgo/internal/pkg/errorsx"
	genericoptions "fastgo/pkg/options"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	go func() {
		// s.srv是一个http服务实例,调用方法开始监听客户端请求
		// http.ErrServerClosed意味着服务器正常关闭
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	// 实现优雅关闭
	// 创建一个os.Singal类型的channel, 用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞主进程, 等待从 quit channel 中接收到信号
	<-quit

	slog.Info("Shutting down server ...")
	// 优雅关停服务
	// 创建上下文对象 ctx, 指定超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 先关闭依赖的服务, 再关闭被依赖的服务
	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Insecure Server forced to shutdown", "err", err)
		return err
	}

	// 正常关闭
	slog.Info("Server exited")
	return nil
}

// NewServer 根据配置创建服务器.
func (cfg *Config) NewServer() (*Server, error) {
	// 创建gin引擎.
	engine := gin.New()

	// 初始化数据库连接
	db, err := cfg.MySQLOptions.NewDB()
	if err != nil {
		return nil, err
	}
	store := store2.NewStore(db)
	cfg.InstallRESTAPI(engine, store)

	//// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	//mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.RequestID()}
	//// Use()函数入参接收一个可变参数, 但mws是一个切片, 切片后加...可以将其解构为多个独立参数
	//engine.Use(mws...)
	//// 注册 404 Handler.
	//engine.NoRoute(func(c *gin.Context) {
	//	core.WriteResponse(c, errorsx.ErrNotFound.WithMessage("Page not found"), nil)
	//})
	//// 注册 /healthz handler.
	//// 请求方法: GET; 请求路径: /healthz; 请求返回: {"status":"ok"}
	//engine.GET("/healthz", func(c *gin.Context) {
	//	core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	//})

	// 创建 HTTP Server 实例.
	// 将cfg配置和http服务器实例注入到新创的结构体中
	httpsrv := &http.Server{Addr: cfg.Addr, Handler: engine}

	return &Server{
		cfg: cfg,
		srv: httpsrv,
	}, nil
}

func (cfg *Config) InstallRESTAPI(engine *gin.Engine, store store2.IStore) {
	// 注册 404 Handler
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, nil, errorsx.ErrNotFound.WithMessage("Page not found"))
	})

	// 注册 /healthz handler.
	engine.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// 创建业务处理器Handler
	handler := handler.NewHandler(biz.NewBiz(store), validation.NewValidator(store))

	// gin.HandlerFunc类型的切片
	// 是用来处理HTTP请求的函数类型, 作用是为路由分组添加中间件.
	authMiddlewares := []gin.HandlerFunc{AuthMiddleware()}

	// 注册 v1 版本 API 路由分组
	v1 := engine.Group("/v1")
	{
		// 用户模块相关路由
		userv1 := v1.Group("/users")
		{
			userv1.POST("", handler.CreateUser) // 创建用户
			// 更新用户信息
			// 删除用户
			// 查询用户详情
			// 查询用户列表
		}
		// 博客模块相关路由
		// 所有以/v1/posts开头的路由都会先经过authMiddlewares里的中间件处理. 只有通过了身份验证中间件的验证, 请求才会被转发到对应的处理函数.
		postv1 := v1.Group("/posts", authMiddlewares...)
		{
			// 创建博客
			// 更新博客
			// 删除博客
			// 查询博客详情
			// 查询博客列表
		}
	}

}

// AuthMiddleware 是一个简单的身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 假设这里进行身份验证逻辑，例如检查请求头中的token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// 这里可以添加更多的验证逻辑，比如验证token的有效性

		// 如果验证通过，继续处理请求
		c.Next()
	}
}
