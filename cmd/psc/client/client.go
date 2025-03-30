package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Subscribe 订阅消息
func Subscribe(addr, token, namespace, topic string) {
	// 验证必填参数
	if namespace == "" || topic == "" || token == "" {
		fmt.Println("错误: 命名空间、主题和令牌都是必填项")
		os.Exit(1)
	}
	var (
		lastID int64 = 0
	)
	for {
		// 构建请求URL
		url := fmt.Sprintf("http://%s/s/%s/%s?lastID=%d", addr, namespace, topic, lastID)
		
		// 创建请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("创建请求失败: %v\n", err)
			continue
		}
		
		// 设置认证头
		req.Header.Set("Authorization", token)
		
		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("连接服务器失败: %v\n", err)
			continue
		}
		defer resp.Body.Close()
		
		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("订阅失败: (%d) %s\n", resp.StatusCode, string(body))
			os.Exit(1)
		}
		
		// 读取 SSE 数据流
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("读取数据失败: %v\n", err)
				break
			}
			
			// 解析 SSE 格式数据
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				data = strings.TrimSpace(data)
				
				// 解析消息
				var msg struct {
					ID        int64           `json:"ID"`
					Data      json.RawMessage `json:"Data"`
					Timestamp int64           `json:"Timestamp"`
				}
				
				if err := json.Unmarshal([]byte(data), &msg); err != nil {
					fmt.Printf("解析消息失败: %v\n", err)
					continue
				}
				lastID = msg.ID
				fmt.Println(data)
			}
		}
	}
}

// Publish 发布消息
func Publish(addr, token, namespace, topic, data string) {
	// 验证必填参数
	if namespace == "" || topic == "" || token == "" || data == "" {
		fmt.Println("错误: 命名空间、主题、令牌和数据都是必填项")
		os.Exit(1)
	}
	fmt.Println("publish", namespace, topic, data)
	// 构建请求URL
	url := fmt.Sprintf("http://%s/p/%s/%s", addr, namespace, topic)
	
	// 创建消息
	msg := struct {
		ID        int64           `json:"ID"`
		Data      json.RawMessage `json:"Data"`
		Timestamp int64           `json:"Timestamp"`
	}{
		ID:        time.Now().UnixNano() / 1000000,            // 毫秒时间戳作为ID
		Data:      json.RawMessage(fmt.Sprintf(`"%s"`, data)), // 简单文本消息
		Timestamp: time.Now().Unix(),
	}
	
	// 序列化消息
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("序列化消息失败: %v\n", err)
		os.Exit(1)
	}
	
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msgBytes))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		os.Exit(1)
	}
	
	// 设置认证头和内容类型
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("连接服务器失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("发布失败: %s (%d)\n", string(body), resp.StatusCode)
		os.Exit(1)
	}
	
	fmt.Printf("消息已发布到 %s/%s\n", namespace, topic)
}
