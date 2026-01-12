package option

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

// 错误码模块
const ModuleOption = 20

// 领域错误定义
var (
	// ErrDatabaseError 数据库错误
	ErrDatabaseError = errcode.Register(errcode.New(
		ModuleOption, 1001,
		"option",
		"error.option.database_error",
		"数据库操作失败",
		http.StatusInternalServerError,
	))

	// ErrNotFound 配置不存在
	ErrNotFound = errcode.Register(errcode.New(
		ModuleOption, 1002,
		"option",
		"error.option.not_found",
		"配置不存在",
		http.StatusNotFound,
	))

	// ErrKeyExists 配置键已存在
	ErrKeyExists = errcode.Register(errcode.New(
		ModuleOption, 1003,
		"option",
		"error.option.key_exists",
		"配置键已存在",
		http.StatusConflict,
	))
)
