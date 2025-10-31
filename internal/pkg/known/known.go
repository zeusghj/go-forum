package known

const (
	// XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 uuid.
	XRequestIDKey = "X-Request-ID"

	// ZhCN 简体中文 - 中国
	ZhCN = "zh-cn"

	// EnUS 英文 - 美国
	EnUS = "en-us"

	// XUserIDKey XUsernameKey 用来定义 Gin 上下文的键，代表请求的所有者.
	XUserIDKey   = "X-UserID"
	XUsernameKey = "X-Username"
	XUserRoleKey = "X-UserRole"
)
