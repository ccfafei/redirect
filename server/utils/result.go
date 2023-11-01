package utils

import (
	"net/http"
	"redirect/model"
	"time"
)

type ResultJson model.ResultJson

// ResultJsonSuccess 返回成功结果
func ResultJsonSuccess() ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Status:  true,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonSuccessWithData 返回成功结果
func ResultJsonSuccessWithData(data interface{}) ResultJson {
	return ResultJson{
		Code:    http.StatusOK,
		Message: "success",
		Result:  data,
		Status:  true,
		Date:    time.Now(),
	}
}

// ResultJsonError 返回错误结果
func ResultJsonError(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusInternalServerError,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonBadRequest 返回错误结果
func ResultJsonBadRequest(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusBadRequest,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}

// ResultJsonUnauthorized 返回错误结果
func ResultJsonUnauthorized(message string) ResultJson {
	return ResultJson{
		Code:    http.StatusUnauthorized,
		Message: message,
		Status:  false,
		Result:  nil,
		Date:    time.Now(),
	}
}
