// 实现上下文相关逻辑
package ruo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 用于返回json格式的数据
type H map[string]interface{}

// 上下文结构体
type Context struct {
	// 原始对象
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求信息
	Path   string
	Method string

	Params map[string]string

	// 响应信息
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int
}

// 初始化上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 处理中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c) // 执行下一个中间件
	}
}

// 获取 post 数据
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取 querystring 的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置返回状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(c.StatusCode)
}

// 设置header
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// 获取header
func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

// 返回格式 string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回格式 json
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// 返回格式 html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	s, _ := c.Params[key]
	return s
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.Json(code, H{"message": err})
}
