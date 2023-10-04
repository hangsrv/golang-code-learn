package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	responseOKFormat       = "HTTP/1.1 200 OK\r\n\r\n%s"
	responseNOTFOUNDFormat = "HTTP/1.1 404 Not Found\r\n\r\n%s"
)

type request struct {
	requestURI  string
	requestBody []byte
}

type response struct {
	c net.Conn
}

type router struct {
	routes map[string]func(*request, *response)
}

func newRouter() *router {
	return &router{
		routes: make(map[string]func(*request, *response)),
	}
}

func (r *router) AddRoute(path string, handler func(*request, *response)) {
	r.routes[path] = handler
}

func (r *router) serveHTTP(req *request, rsp *response) {
	handler, ok := r.routes[req.requestURI]
	if ok {
		handler(req, rsp)
	} else {
		notFoundHandler(req, rsp)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	defer listener.Close()
	log.Println("Server listening on ", ":8080")
	router := newRouter()
	router.AddRoute("/", rootHandler)
	router.AddRoute("/hello", pingHandler)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn, router)
	}
}

func handleRequest(conn net.Conn, router *router) {
	defer conn.Close()
	req := &request{}
	rsp := &response{c: conn}

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading request line:", err)
		return
	}

	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		log.Println("Invalid request format")
		return
	}

	requestMethod := parts[0]
	requestURI := parts[1]
	httpVersion := parts[2]

	log.Printf("Request Method: %s\n", requestMethod)
	log.Printf("Request URI: %s\n", requestURI)
	log.Printf("HTTP Version: %s\n", httpVersion)

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headerName := strings.TrimSpace(parts[0])
			headerValue := strings.TrimSpace(parts[1])
			headers[headerName] = headerValue
		}
	}

	// 读取请求体（如果有）
	contentLength, ok := headers["Content-Length"]
	if ok {
		bodySize := 0
		for bodySize < len(contentLength) {
			// 读取请求体数据
			bodyData := make([]byte, 1024)
			n, err := reader.Read(bodyData)
			if err != nil || n == 0 {
				break
			}

			// 处理请求体数据
			bodySize += n
			req.requestBody = append(req.requestBody, bodyData[:n]...)
		}
	}
	log.Println("req ", req)
	
	req.requestURI = requestURI

	router.serveHTTP(req, rsp)
}

func rootHandler(req *request, rsp *response) {
	response := fmt.Sprintf(responseOKFormat, "hello redis")
	rsp.c.Write([]byte(response))
	rsp.c.Close()
}

func pingHandler(req *request, rsp *response) {
	response := fmt.Sprintf(responseOKFormat, "ping!")
	rsp.c.Write([]byte(response))
	rsp.c.Close()
}

func notFoundHandler(req *request, rsp *response) {
	response := fmt.Sprintf(responseNOTFOUNDFormat, "Page not found")
	rsp.c.Write([]byte(response))
	rsp.c.Close()
}
