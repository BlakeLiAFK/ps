package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// getInt64Query 获取整数查询参数
func getInt64Query(r *http.Request, key string, defaultValue int64) int64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// sendMessageSSE 将消息发送为SSE格式
func sendMessageSSE(w http.ResponseWriter, flusher http.Flusher, msg *Message) {
	fmt.Println("send data", msg)
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}
