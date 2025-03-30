package main

import (
	"github.com/BlakeLiAFK/ps/cmd/psc/client"
	"github.com/BlakeLiAFK/ps/cmd/psc/server"
	"github.com/spf13/cobra"
)

// 可以在编译时通过-ldflags注入的变量
var (
	DefaultAddr string = "localhost:8080"
)

var (
	// 全局选项
	addr         string
	token        string
	isServerMode bool

	// client专用
	namespace string
	topic     string
	data      string
)

func main() {
	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "psc",
		Short: "PubSub客户端工具",
		Run: func(cmd *cobra.Command, args []string) {
			// 默认行为 - 根据是否有-s标志决定模式
			if isServerMode {
				server.Run(addr, token)
			} else {
				// 根据是否提供data参数决定是pub还是sub
				if data != "" {
					client.Publish(addr, token, namespace, topic, data)
				} else {
					client.Subscribe(addr, token, namespace, topic)
				}
			}
		},
	}

	// 设置全局选项
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", DefaultAddr, "服务器地址")
	rootCmd.PersistentFlags().StringVarP(&token, "key", "k", "", "认证令牌 (必填)")
	rootCmd.PersistentFlags().BoolVarP(&isServerMode, "server", "s", false, "以服务器模式运行")

	// 设置客户端选项
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "命名空间")
	rootCmd.PersistentFlags().StringVarP(&topic, "topic", "t", "", "主题")
	rootCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "要发布的数据")

	// 执行命令
	rootCmd.Execute()
}
