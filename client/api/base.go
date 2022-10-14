package api

import "encoding/json"

func MarshalJsonData(str interface{}) string {
	marshal, err := json.Marshal(str)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func UnmarshalJsonData(byte []byte, target interface{}) {
	err := json.Unmarshal(byte, target)
	if err != nil {
		return
	}
}
