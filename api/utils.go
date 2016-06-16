package api

import (
	"encoding/json"
	"net/http"
)

// 解析http.request中body参数到实体
func DecoderRequest(req *http.Request, struzt interface{}) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(struzt)
}
