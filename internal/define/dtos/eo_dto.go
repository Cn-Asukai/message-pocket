package dtos

// EOEventRequest EdgeOne 页面事件请求
type EOEventRequest struct {
	EventType    string `json:"eventType"`
	AppID        string `json:"appId"`
	ProjectID    string `json:"projectId"`
	DeploymentID string `json:"deploymentId"`
	ProjectName  string `json:"projectName"`
	RepoBranch   string `json:"repoBranch"`
	Timestamp    string `json:"timestamp"`
}
