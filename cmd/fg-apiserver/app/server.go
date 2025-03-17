package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 该方法创建一个*cobra.Command对象, 用于启动应用程序
func NewFastGOCommand() *cobra.Command {
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
			fmt.Println("hello fastgo!")
			return nil
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./fg-apiserver param1 param2
		Args: cobra.NoArgs,
	}
	return cmd
}
