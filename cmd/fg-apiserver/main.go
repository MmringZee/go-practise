package main

import (
	"fastgo/cmd/fg-apiserver/app"
	_ "go.uber.org/automaxprocs"
	"os"
)

func main() {
	// 创建Go项目
	command := app.NewFastGOCommand()

	// 执行命令并处理错误
	if err := command.Execute(); err != nil {
		// err != nil 意味着发生了异常
		// 返回退出码，可以使其他程序（例如 bash 脚本）根据退出码来判断服务运行状态
		os.Exit(1)
	}
}
