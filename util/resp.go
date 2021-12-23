package util

import (
	"encoding/json"
	"io"
)

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// JSONBytes :
func (m *RespMsg) JsonBytes() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//WriteTo: WriteTo writer
func (m *RespMsg) WriteTo(w io.Writer) (n int64, err error) {
	data, err := m.JsonBytes()
	if err != nil {
		return 0, err
	}
	writeBytesLen, err := w.Write(data)

	return int64(writeBytesLen), err
}
