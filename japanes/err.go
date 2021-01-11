package main

import (
	"github.com/iancoleman/strcase"
)

// MyErr MyErr
type MyErr struct {
	Code         int64
	Msg          string
	ExtraInfo    map[string]interface{}
	OriginErrors []string
}

// Error Code 常數定義
const (
	// 部門定義代碼
	DepartmentCode = 12
	// 服務定義代碼
	ServeiceCode = 16

	SystemErrorCode = 9999
)

// API Error Code 常數定義 自動增加
const (
	APINotExistErrorCode = int64(iota) + 1
	RequsetParamsErrorCode
	GameTypeNotSupportErrorCode
	GameNotExistErrorCode
	PriceNotExistErrorCode
	PriceLevelNotExistErrorCode
	AlreadyInOrderGenerateQueueErrorCode
	NoOrderErrorCode
	MethodNotAllowedErrorCode
)

func formatErrorCode(code int64) int64 {
	return DepartmentCode*10000000 + ServeiceCode*10000 + code
}

// NewMyError NewMyError
func NewMyError(code int64, msg string) *MyErr {
	return &MyErr{
		Code:         formatErrorCode(SystemErrorCode),
		Msg:          strcase.ToScreamingSnake(msg),
		ExtraInfo:    make(map[string]interface{}),
		OriginErrors: make([]string, 0),
	}
}

// SystemError SystemError
func SystemError(originErr error) *MyErr {
	e := NewMyError(999, "SystemError")
	if originErr != nil {
		e.AddOriginError(originErr)
	}

	return e
}

// AddOriginError AddOriginError
func (e *MyErr) AddOriginError(err error) {
	e.OriginErrors = append(e.OriginErrors, err.Error())
}

// AddExtraInfo AddExtraInfo
func (e *MyErr) AddExtraInfo(k string, v interface{}) {
	e.ExtraInfo[k] = v
}

// RequsetParamsError RequsetParamsError
func RequsetParamsError(k string, v interface{}) *MyErr {
	e := NewMyError(RequsetParamsErrorCode, "RequsetParamsError")
	e.AddExtraInfo(k, v)

	return e
}

// RequsetParamsError RequsetParamsError
func NoOrdersError(gameId int64, wordIdx int8) *MyErr {
	e := NewMyError(NoOrderErrorCode, "NoOrderErrorCode")
	e.AddExtraInfo("gameId", gameId)
	e.AddExtraInfo("wordIdx", wordIdx)

	return e
}
