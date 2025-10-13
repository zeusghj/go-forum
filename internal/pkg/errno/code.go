// Copyright 2022 zeusghj(郭洪军) <zeusghj@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/zeusghj/scaffold.

package errno

import (
	"go-forum/internal/pkg/known"
	"net/http"

	"github.com/spf13/viper"
)

var (
	// OK 代表请求成功
	OK = &Errno{HTTP: http.StatusOK, Code: 200, Message: ""}

	// InternalServerError 表示所有未知的服务器端错误.
	InternalServerError = &Errno{HTTP: http.StatusInternalServerError, Code: 500, Message: "Internal server error."}

	// ErrPageNotFound 表示路由不匹配错误.
	ErrPageNotFound = &Errno{HTTP: http.StatusNotFound, Code: ResourcePageNotFound, Message: Text(ResourcePageNotFound)}

	// ErrUserAlreadyExist 用户已存在
	ErrUserAlreadyExist = &Errno{HTTP: http.StatusOK, Code: UserAlreadyExistError, Message: Text(UserAlreadyExistError)}

	// ErrBind 表示参数绑定错误.
	ErrBind = &Errno{HTTP: http.StatusOK, Code: BindError, Message: Text(BindError)}

	// ErrInvalidParameter 表示所有验证失败的错误.
	ErrInvalidParameter = &Errno{HTTP: http.StatusOK, Code: InvalidParameterError, Message: Text(InvalidParameterError)}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{HTTP: http.StatusUnauthorized, Code: SignTokenError, Message: Text(SignTokenError)}

	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{HTTP: http.StatusUnauthorized, Code: AuthFailure, Message: Text(AuthFailure)}

	// ErrPasswordIncorrect 密码不正确
	ErrPasswordIncorrect = &Errno{HTTP: http.StatusOK, Code: PasswordInCorrect, Message: Text(PasswordInCorrect)}

	// ErrUserNotFound 用户不存在
	ErrUserNotFound = &Errno{HTTP: http.StatusOK, Code: UserNotFound, Message: Text(UserNotFound)}

	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{HTTP: http.StatusUnauthorized, Code: Unauthorized, Message: Text(Unauthorized)}

	// ErrPostNotFound 表示请求的文章不存在.
	ErrPostNotFound = &Errno{HTTP: http.StatusOK, Code: PostNotFoundError, Message: Text(PostNotFoundError)}
)

const (
	ServerError           = 10101
	TooManyRequests       = 10102
	ParamBindError        = 10103
	AuthorizationError    = 10104
	UrlSignError          = 10105
	CacheSetError         = 10106
	CacheGetError         = 10107
	CacheDelError         = 10108
	CacheNotExist         = 10109
	ResubmitError         = 10110
	HashIdsEncodeError    = 10111
	HashIdsDecodeError    = 10112
	RBACError             = 10113
	RedisConnectError     = 10114
	MySQLConnectError     = 10115
	WriteConfigError      = 10116
	SendEmailError        = 10117
	MySQLExecError        = 10118
	GoVersionError        = 10119
	SocketConnectError    = 10120
	SocketSendError       = 10121
	ResourcePageNotFound  = 10122
	BindError             = 10123
	InvalidParameterError = 10124
	AuthFailure           = 10225
	SignTokenError        = 10226
	PasswordInCorrect     = 10227
	UserNotFound          = 10228
	Unauthorized          = 10229

	AuthorizedCreateError    = 20101
	AuthorizedListError      = 20102
	AuthorizedDeleteError    = 20103
	AuthorizedUpdateError    = 20104
	AuthorizedDetailError    = 20105
	AuthorizedCreateAPIError = 20106
	AuthorizedListAPIError   = 20107
	AuthorizedDeleteAPIError = 20108

	AdminCreateError             = 20201
	AdminListError               = 20202
	AdminDeleteError             = 20203
	AdminUpdateError             = 20204
	AdminResetPasswordError      = 20205
	AdminLoginError              = 20206
	AdminLogOutError             = 20207
	AdminModifyPasswordError     = 20208
	AdminModifyPersonalInfoError = 20209
	AdminMenuListError           = 20210
	AdminMenuCreateError         = 20211
	AdminOfflineError            = 20212
	AdminDetailError             = 20213

	MenuCreateError       = 20301
	MenuUpdateError       = 20302
	MenuListError         = 20303
	MenuDeleteError       = 20304
	MenuDetailError       = 20305
	MenuCreateActionError = 20306
	MenuListActionError   = 20307
	MenuDeleteActionError = 20308

	CronCreateError  = 20401
	CronUpdateError  = 20402
	CronListError    = 20403
	CronDetailError  = 20404
	CronExecuteError = 20405

	UserAlreadyExistError = 20501
	PostNotFoundError     = 20502
)

func Text(code int) string {
	lang := viper.GetString("language.local")

	if lang == known.ZhCN {
		return zhCNText[code]
	}

	if lang == known.EnUS {
		return enUSText[code]
	}

	return zhCNText[code]
}
