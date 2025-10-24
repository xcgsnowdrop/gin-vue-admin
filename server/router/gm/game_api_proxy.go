package gm

import (
	"fmt"
	"io"
	"net/http"

	"gmserver/global"
	"gmserver/middleware"
	"gmserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type GameApiProxyRouter struct{}

func (s *GameApiProxyRouter) InitGameApiProxyRouter(Router *gin.RouterGroup) {
	// 创建一个通用的代理路由，匹配所有路径
	// 注意：PrivateGroup已经应用了JWTAuth()和CasbinHandler()中间件
	// 这里只需要添加操作记录中间件
	Router.Group("gm").
		Use(middleware.OperationRecord()).
		Any("/*path", s.gameApiProxy)
}

// gameApiProxy 游戏API代理实现
func (s *GameApiProxyRouter) gameApiProxy(c *gin.Context) {
	// 获取原始路径
	originalPath := c.Request.URL.Path

	// 构建目标URL
	targetURL := fmt.Sprintf("http://%s:%d%s",
		global.GVA_CONFIG.GameAPI.Host,
		global.GVA_CONFIG.GameAPI.Port,
		originalPath)

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取请求体失败"})
		return
	}

	// 创建HTTP客户端
	client := resty.New()
	client.SetTimeout(global.GVA_CONFIG.GameAPI.Timeout)
	client.SetRetryCount(global.GVA_CONFIG.GameAPI.RetryCount)

	// 复制请求头
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// 转换查询参数
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// 发送请求
	var resp *resty.Response
	switch c.Request.Method {
	case "GET":
		resp, err = client.R().
			SetHeaders(headers).
			SetQueryParams(queryParams).
			Get(targetURL)
	case "POST":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(body).
			Post(targetURL)
	case "PUT":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(body).
			Put(targetURL)
	case "DELETE":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(body).
			Delete(targetURL)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "不支持的请求方法"})
		return
	}

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("游戏API请求失败: %v, URL: %s", err, targetURL))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "游戏API请求失败"})
		return
	}

	// 记录操作日志
	userID := utils.GetUserID(c)
	global.GVA_LOG.Info(fmt.Sprintf("用户 %d 请求游戏API: %s %s -> %s (状态码: %d)",
		userID, c.Request.Method, originalPath, targetURL, resp.StatusCode()))

	// 复制响应头
	for key, values := range resp.Header() {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 返回响应
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
