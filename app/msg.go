package app

type ClientControlResponse struct {
	GeneraResponse
}

type GeneraResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type ClientControlRequest struct {
	Type    int64  `json:"type"`
	AppName string `json:"app_name"`
	Passwd  string `json:"passwd"`
}
