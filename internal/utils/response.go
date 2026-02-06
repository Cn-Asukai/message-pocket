package utils

import "github.com/samber/lo"

func NewJsonResponseWithoutData(code int, message string) map[string]any {
	return NewJsonResponse(code, message, nil)
}

func NewJsonResponse(code int, message string, data any) map[string]any {
	resp := map[string]any{
		"code":    code,
		"message": message,
	}
	if !lo.IsEmpty(data) {
		resp["data"] = data
	}

	return resp
}
