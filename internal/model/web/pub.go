package web

// GetCaptchaResp 定义了获取验证码响应结构体。
type GetCaptchaResp struct {
	Id      string `json:"id"`
	Content string `json:"b64s"`
}

// SendCaptchaReq 定义了发送验证码请求结构体。
type SendCaptchaReq struct {
	Phone        string `json:"phone" form:"phone" binding:"required"`
	ImgCaptcha   string `json:"img_captcha" form:"img_captcha" binding:"required"`
	ImgCaptchaID string `json:"img_captcha_id" form:"img_captcha_id" binding:"required"`
}

// LoginReq 定义了登录请求结构体。
type LoginReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// LoginResp 定义了登录响应结构体。
type LoginResp struct {
	Token string `json:"token"`
}

// RegisterReq 定义了注册请求结构体。
type RegisterReq struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required"`
	Oauth     string `json:"oauth" form:"oauth" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
}

// RegisterResp 定义了注册响应结构体。
type RegisterResp struct {
	UserId   int64  `json:"id"`
	Username string `json:"username"`
}
