package utils

import (
	"encoding/json"
	"fmt"

	"gmserver/global"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// GameAPIResponse 游戏API响应结构
type GameAPIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// CallGameAPI 调用游戏API的通用方法
// @param path 请求路径，例如 "/email/system/send"
// @param method 请求方法，例如 "POST", "GET" 等
// @param payload 请求体数据（如果是POST/PUT等需要body的请求）
// @return *GameAPIResponse 游戏API响应
// @return error 错误信息
func CallGameAPI(path string, method string, payload interface{}) (*GameAPIResponse, error) {
	// 构建完整的URL
	url := fmt.Sprintf("http://%s:%d%s",
		global.GVA_CONFIG.GameAPI.Host,
		global.GVA_CONFIG.GameAPI.Port,
		path)

	// 创建HTTP客户端
	client := resty.New()
	client.SetTimeout(global.GVA_CONFIG.GameAPI.Timeout)
	client.SetRetryCount(global.GVA_CONFIG.GameAPI.RetryCount)

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// 添加Authorization请求头（如果配置中存在）
	if global.GVA_CONFIG.GameAPI.Authorization != "" {
		headers["Authorization"] = global.GVA_CONFIG.GameAPI.Authorization
	}

	var resp *resty.Response
	var err error

	// 根据请求方法发送请求
	switch method {
	case "GET":
		resp, err = client.R().
			SetHeaders(headers).
			Get(url)
	case "POST":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(payload).
			Post(url)
	case "PUT":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(payload).
			Put(url)
	case "DELETE":
		resp, err = client.R().
			SetHeaders(headers).
			SetBody(payload).
			Delete(url)
	default:
		return nil, fmt.Errorf("不支持的请求方法: %s", method)
	}

	if err != nil {
		global.GVA_LOG.Error("游戏API请求失败",
			zap.String("url", url),
			zap.String("method", method),
			zap.Error(err))
		return nil, fmt.Errorf("游戏API请求失败: %v", err)
	}

	// 解析响应
	var apiResp GameAPIResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		global.GVA_LOG.Error("解析游戏API响应失败",
			zap.String("url", url),
			zap.String("response", string(resp.Body())),
			zap.Error(err))
		return nil, fmt.Errorf("解析游戏API响应失败: %v", err)
	}

	// 记录日志
	global.GVA_LOG.Info("游戏API请求成功",
		zap.String("url", url),
		zap.String("method", method),
		zap.Int("status_code", resp.StatusCode()),
		zap.Int("response_code", apiResp.Code),
		zap.String("response_msg", apiResp.Msg))

	return &apiResp, nil
}

// CallGameAPIPOST 调用游戏API的POST方法（便捷方法）
func CallGameAPIPOST(path string, payload interface{}) (*GameAPIResponse, error) {
	return CallGameAPI(path, "POST", payload)
}

// CallGameAPIGET 调用游戏API的GET方法（便捷方法）
func CallGameAPIGET(path string) (*GameAPIResponse, error) {
	return CallGameAPI(path, "GET", nil)
}
