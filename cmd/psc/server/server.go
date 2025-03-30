package server

import (
	"fmt"
	"os"
	
	"github.com/BlakeLiAFK/ps"
)

// Run 运行服务器模式
func Run(addr string, token string) {
	if token == "" {
		fmt.Println("错误: 认证令牌是必填项")
		os.Exit(1)
	}
	
	fmt.Printf("启动服务器于 %s, 使用令牌: %s\n", addr, token)
	
	// 创建认证器和服务器
	auth := ps.NewAuth(token)
	ps.NewServer(auth).Run(addr)
}
