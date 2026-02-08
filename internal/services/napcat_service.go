package services

import (
	"context"
	"errors"
	"fmt"
	"message-pocket/internal/config"

	"github.com/samber/do/v2"
	"resty.dev/v3"
)

// Response NapCat API 响应结构
type Response[T any] struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Wording string `json:"wording"`
	Stream  string `json:"stream"`
}

// NapCatService NapCat 服务
type NapCatService struct {
	token  string
	apiURL string
}

// NewNapCatService 创建 NapCat 服务实例
func NewNapCatService(cfg *config.Config) *NapCatService {
	return &NapCatService{
		token:  cfg.NapCatConfig.Token,
		apiURL: cfg.NapCatConfig.URL,
	}
}

func ProvideNapCatService(i do.Injector) (*NapCatService, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return NewNapCatService(cfg), nil
}

// getURL 构建完整的 API URL
func (s *NapCatService) getURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", s.apiURL, endpoint)
}

// post 发送 POST 请求到 NapCat API
func (s *NapCatService) post(ctx context.Context, endpoint string, body any) (*any, error) {
	url := s.getURL(endpoint)

	client := resty.New()
	responseData := Response[any]{}
	response, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.token)).
		SetBody(body).
		SetResult(&responseData).
		Post(url)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, errors.New(fmt.Sprintf("status code: %d, body: %s", response.StatusCode(), responseData.Message))
	}

	// 检查 API 状态
	if responseData.Status != "ok" {
		return nil, fmt.Errorf("NapCat API error: %s (retcode: %d)", responseData.Message, responseData.Retcode)
	}

	return &responseData.Data, nil
}

// SendGroupMessage 发送群消息
func (s *NapCatService) SendGroupMessage(ctx context.Context, groupID, message string) error {
	type SendGroupMsgRequest struct {
		GroupID string `json:"group_id"`
		Message string `json:"message"`
	}

	type SendGroupMsgResponse struct {
		MessageID int64 `json:"message_id"`
	}

	req := SendGroupMsgRequest{
		GroupID: groupID,
		Message: message,
	}

	_, err := s.post(ctx, "send_group_msg", req)
	return err
}
