package app

import (
	"fastgo/cmd/fg-apiserver/app/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"os"
)

var configFile string // 配置文件路径

// 该方法创建一个*cobra.Command对象, 用于启动应用程序
func NewFastGOCommand() *cobra.Command {

	// 创建默认的应用命令行选项
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		// 指定命令的名称, 改名字会出现在帮助信息中
		Use: "fg-apiserver",
		// 命令描述
		Short: "A very lightweight full go project",
		Long:  "A very lightweight full go project, designed to help beginners quickly learn Go project development.",
		// 命令出错时，不打印帮助信息。设置为true可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 调用cmd.Execute()时执行的Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},

		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./fg-apiserver param1 param2
		Args: cobra.NoArgs,
	}

	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the fg-apiserver configuration file.")

	return cmd
}

// run 是主运行逻辑，负责初始化日志、解析配置、校验选项并启动服务器.
func run(opts *options.ServerOptions) error {

	// 初始化 slog
	initLog()

	// 将 viper 中的配置解析到 opts.
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	// 校验命令行选项
	if err := opts.Validate(); err != nil {
		return err
	}

	// 获取应用配置.
	// 将命令行选项和应用配置分开，可以更加灵活的处理 2 种不同类型的配置.
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// 创建服务器实例.
	server, err := cfg.NewServer()
	if err != nil {
		return err
	}

	return server.Run()
}

// initLog 初始化全局日志实例
func initLog() {
	// 获取日志配置
	// 通过viper获取配置文件中的键值
	format := viper.GetString("log.format") // 日志格式, 支持: json、text
	level := viper.GetString("log.level")   // 日志级别，支持：debug, info, warn, error
	output := viper.GetString("log.output") // 日志输出路径，支持：标准输出stdout和文件

	// 转换日志级别
	var slevel slog.Level
	switch level {
	case "debug":
		slevel = slog.LevelDebug
	case "info":
		slevel = slog.LevelInfo
	case "warn":
		slevel = slog.LevelWarn
	case "error":
		slevel = slog.LevelWarn
	default:
		slevel = slog.LevelInfo
	}

	// slog/log在创建Handler时提供了一个配置
	opts := &slog.HandlerOptions{Level: slevel}

	// 转换日志输出格式
	var w io.Writer
	var err error
	switch output {
	case "":
		w = os.Stdout
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 转换日志格式
	if err != nil {
		return
	}
	var handler slog.Handler
	switch format {
	// 以json格式输出
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	// 以key=value
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)
	}

	// 设置全局的日志实例为自定义的日志实例
	slog.SetDefault(slog.New(handler))
}
