package dto

import "net/http"

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func NewResp[T any](code int, msg string, data T) Resp[T] {
	return Resp[T]{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewRespWithSuccessData[T any](data T) Resp[T] {
	return NewResp(http.StatusOK, "success", data)
}
