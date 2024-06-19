package utils

import "encoding/json"

func ToJsonString(v any) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

func ToJsonByte(v any) []byte {
	buf, _ := json.Marshal(v)
	return buf
}
