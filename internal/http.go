package internal

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

type httpServer struct {
	server *http.Server
	mux    *http.ServeMux
	Ctx    *Context
}

func newHttpServer(ctx *Context) *httpServer {
	mux := http.NewServeMux()
	s := &httpServer{
		Ctx: ctx,
		server: &http.Server{
			Handler: mux,
		},
		mux: mux,
	}
	s.initRoute()
	return s
}

func (s *httpServer) initRoute() {
	s.mux.HandleFunc("/s/", s.handleSub)
	s.mux.HandleFunc("/p/", s.handlePub)
}

func (s *httpServer) handleSub(w http.ResponseWriter, r *http.Request) {
	// 只支持GET请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 验证订阅合法性 bearer token
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// 解析请求路径，格式为 /s/{namespace}/{topic}
	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/s/") {
		http.Error(w, "Invalid path format, expected /s/{namespace}/{topic}", http.StatusBadRequest)
		return
	}
	urlPath = urlPath[len("/s/"):]
	parts := strings.Split(urlPath, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path format, expected /s/{namespace}/{topic}", http.StatusBadRequest)
		return
	}
	namespace := parts[0]
	topic := parts[1]
	// 验证token
	verifyToken, err := s.Ctx.Auth.VerifyToken(token, namespace, topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !verifyToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ns := s.Ctx.PS.GetOrNewNamespace(namespace)
	if ns == nil { // 命名空间进黑名单才会创建失败
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	topicCtx := ns.GetOrNewTopic(topic)
	// 设置SSE头部
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许跨域访问
	// 返回200 OK
	w.WriteHeader(http.StatusOK)
	// 检查flusher支持
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	
	// 获取lastID参数
	lastID := getInt64Query(r, "lastID", 0)
	if lastID != 0 {
		// 获取lastID之后的消息
		messages := topicCtx.GetMessagesFromID(lastID)
		for _, msg := range messages {
			sendMessageSSE(w, flusher, msg)
		}
	}
	// 创建订阅者
	ch := make(chan *Message, 10) // 缓冲区大小设为10
	topicCtx.AddSubscriber(ch)
	// 保持连接
	for {
		select {
		case msg := <-ch:
			// 发送SSE消息
			sendMessageSSE(w, flusher, msg)
		}
	}
}

func (s *httpServer) handlePub(w http.ResponseWriter, r *http.Request) {
	// 只支持POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 验证订阅合法性 bearer token
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// 解析请求路径，格式为 /p/{namespace}/{topic}
	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/p/") {
		http.Error(w, "Invalid path format, expected /p/{namespace}/{topic}", http.StatusBadRequest)
		return
	}
	urlPath = urlPath[len("/p/"):]
	parts := strings.Split(urlPath, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path format, expected /p/{namespace}/{topic}", http.StatusBadRequest)
		return
	}
	namespace := parts[0]
	topic := parts[1]
	// 验证token
	verifyToken, err := s.Ctx.Auth.VerifyToken(token, namespace, topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !verifyToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ns := s.Ctx.PS.GetOrNewNamespace(namespace)
	if ns == nil { // 命名空间进黑名单才会创建失败
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	// 获取消息体
	var msg Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 发布消息
	ns.GetOrNewTopic(topic).Publish(&msg)
	// 返回200 OK
	w.WriteHeader(http.StatusOK)
	// 成功
	_, _ = w.Write([]byte("OK"))
}

func (s *httpServer) Run(addr string) error {
	// 启动 http server
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 确保设置了Handler
	s.server.Handler = s.mux
	return s.server.Serve(ln)
}
