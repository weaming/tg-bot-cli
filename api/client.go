package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string, proxyURL string) (*Client, error) {
	transport := &http.Transport{}

	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("代理地址格式错误: %w", err)
		}
		transport.Proxy = http.ProxyURL(parsed)
	} else {
		// 自动读取 HTTP_PROXY / HTTPS_PROXY 环境变量
		transport.Proxy = http.ProxyFromEnvironment
	}

	return &Client{
		token:      token,
		httpClient: &http.Client{Transport: transport},
	}, nil
}

func (c *Client) apiURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", c.token, method)
}

// doJSON 发送 JSON POST 请求并返回原始响应体
func (c *Client) doJSON(method string, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL(method), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// doMultipart 上传文件并附带额外表单字段
func (c *Client) doMultipart(apiMethod string, fields map[string]string, fileField string, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for key, value := range fields {
		if value == "" {
			continue
		}
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.apiURL(apiMethod), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// parseResponse 解析通用响应，出错时返回 error
func parseResponse[T any](body []byte) (*T, error) {
	var result Response[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	if !result.OK {
		return nil, fmt.Errorf("API 错误 %d: %s", result.ErrorCode, result.Description)
	}
	return &result.Result, nil
}

func (c *Client) GetMe() (*User, error) {
	body, err := c.doJSON("getMe", struct{}{})
	if err != nil {
		return nil, err
	}
	return parseResponse[User](body)
}

func (c *Client) GetChat(chatID string) (*Chat, error) {
	body, err := c.doJSON("getChat", map[string]string{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	return parseResponse[Chat](body)
}

func (c *Client) SendMessage(params SendMessageParams) (*Message, error) {
	body, err := c.doJSON("sendMessage", params)
	if err != nil {
		return nil, err
	}
	return parseResponse[Message](body)
}

// detectMediaMethod 根据文件扩展名自动选择 API 方法和表单字段名
func detectMediaMethod(path string) (apiMethod, fileField string) {
	ext := strings.ToLower(filepath.Ext(path))

	mimeType := mime.TypeByExtension(ext)

	switch {
	case strings.HasPrefix(mimeType, "image/gif"):
		return "sendAnimation", "animation"
	case strings.HasPrefix(mimeType, "image/"):
		return "sendPhoto", "photo"
	case strings.HasPrefix(mimeType, "video/"):
		return "sendVideo", "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "sendAudio", "audio"
	default:
		return "sendDocument", "document"
	}
}

func (c *Client) SendMedia(params SendMediaParams) (*Message, error) {
	apiMethod, fileField := detectMediaMethod(params.FilePath)

	fields := map[string]string{
		"chat_id":    params.ChatID,
		"caption":    params.Caption,
		"parse_mode": params.ParseMode,
	}
	if params.MessageThreadID > 0 {
		fields["message_thread_id"] = fmt.Sprintf("%d", params.MessageThreadID)
	}
	if params.ReplyToMessageID > 0 {
		fields["reply_to_message_id"] = fmt.Sprintf("%d", params.ReplyToMessageID)
	}
	if params.DisableNotification {
		fields["disable_notification"] = "true"
	}
	if params.ProtectContent {
		fields["protect_content"] = "true"
	}
	if params.ReplyMarkup != nil {
		markupJSON, err := json.Marshal(params.ReplyMarkup)
		if err != nil {
			return nil, err
		}
		fields["reply_markup"] = string(markupJSON)
	}

	body, err := c.doMultipart(apiMethod, fields, fileField, params.FilePath)
	if err != nil {
		return nil, err
	}
	return parseResponse[Message](body)
}

func (c *Client) EditMessageText(params EditMessageTextParams) (*Message, error) {
	body, err := c.doJSON("editMessageText", params)
	if err != nil {
		return nil, err
	}
	return parseResponse[Message](body)
}

func (c *Client) DeleteMessage(params DeleteMessageParams) (bool, error) {
	body, err := c.doJSON("deleteMessage", params)
	if err != nil {
		return false, err
	}
	result, err := parseResponse[bool](body)
	if err != nil {
		return false, err
	}
	return *result, nil
}

func (c *Client) ForwardMessage(params ForwardMessageParams) (*Message, error) {
	body, err := c.doJSON("forwardMessage", params)
	if err != nil {
		return nil, err
	}
	return parseResponse[Message](body)
}

func (c *Client) CopyMessage(params CopyMessageParams) (*Message, error) {
	body, err := c.doJSON("copyMessage", params)
	if err != nil {
		return nil, err
	}
	return parseResponse[Message](body)
}

func (c *Client) PinChatMessage(params PinChatMessageParams) (bool, error) {
	body, err := c.doJSON("pinChatMessage", params)
	if err != nil {
		return false, err
	}
	result, err := parseResponse[bool](body)
	if err != nil {
		return false, err
	}
	return *result, nil
}

func (c *Client) UnpinChatMessage(params UnpinChatMessageParams) (bool, error) {
	body, err := c.doJSON("unpinChatMessage", params)
	if err != nil {
		return false, err
	}
	result, err := parseResponse[bool](body)
	if err != nil {
		return false, err
	}
	return *result, nil
}
