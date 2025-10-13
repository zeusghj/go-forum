package errno

import (
	"fmt"
	"net/http"
)

// Errno 定义了 项目 使用的错误类型.
type Errno struct {
	HTTP    int
	Code    int
	Message string
}

// Error 实现 error 接口中的 `Error` 方法.
func (err *Errno) Error() string {
	return err.Message
}

// SetMessage 设置 Errno 类型错误中的 Message 字段.
func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// Decode 尝试从 err 中解析出业务错误码和错误信息.
func Decode(err error) (int, int, string) {
	if err == nil {
		return http.StatusOK, 200, ""
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	// 默认返回未知错误码和错误信息. 该错误代表服务端出错
	return http.StatusInternalServerError, 500, "Internal Server Error"
}
