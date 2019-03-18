package block

import (
	"encoding/json"
	"log"
)

//将json字符串转换成数组
func Json2Array(jsonStr string) []string {
	var result []string
	if err := json.Unmarshal([]byte(jsonStr),&result); err != nil {
		log.Panic(err)
	}

	return result
}
