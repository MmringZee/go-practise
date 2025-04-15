package options

import (
	"fastgo/internal/apiserver"
	genericoptions "fastgo/pkg/options"
	"fmt"
	"net"
	"strconv"
	"time"
)

type ServerOptions struct {
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	Addr         string                       `json:"addr" mapstructure:"addr"`
	// JWTKey 定义 JWT 密钥.
	JWTKey string `json:"jwt-key" mapsturcture:"jwt-key"`
	// Expiration 定义 JWT token 的过期时间.
	Expiration time.Duration `json:"expiration" mapsturcture:"expiration"`
}

// NewServerOptions 创建带有默认值的 ServerOptions 实例.
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: genericoptions.NewMySQLOptions(),
		Addr:         "0.0.0.0:6666",
	}
}

// Validate 校验 ServerOptions 中的选项是否合法.
func (o *ServerOptions) Validate() error {
	// 校验服务器地址
	if o.Addr == "" {
		return fmt.Errorf("server address cannot be empty")
	}

	// 检查地址格式是否为host:port
	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invalid server address format '%s': %w", o.Addr, err)
	}

	// 验证端口是否为数字且在有效范围内
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid server port: %s", portStr)
	}

	// 校验MySQL配置
	if err := o.MySQLOptions.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		MySQLOptions: o.MySQLOptions,
		Addr:         o.Addr,
	}, nil
}
