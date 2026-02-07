package logic

// GetMessageTypeLabel 根据事件类型获取中文标签
func GetMessageTypeLabel(eventType string) string {
	switch eventType {
	case "deployment.created":
		return "开始部署"
	case "deployment.succeeded":
		return "部署成功"
	case "deployment.failed":
		return "部署失败"
	case "deployment.cancelled":
		return "部署取消"
	case "deployment.rollback":
		return "部署回滚"
	case "deployment.in_progress":
		return "部署进行中"
	case "build.started":
		return "构建开始"
	case "build.succeeded":
		return "构建成功"
	case "build.failed":
		return "构建失败"
	case "project.created":
		return "项目创建"
	case "project.updated":
		return "项目更新"
	case "project.deleted":
		return "项目删除"
	default:
		return eventType
	}
}
