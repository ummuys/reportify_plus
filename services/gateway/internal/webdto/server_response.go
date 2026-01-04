package webdto

type BaseResponse struct {
	Message string `json:"msg"`
}

type ErrResponse struct {
	Error string `json:"err"`
}
