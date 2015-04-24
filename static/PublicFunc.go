package static

import (
	"fmt"
)

//异常处理方法
func DealWithError(info string, err error) bool {
	if err != nil {
		fmt.Println(info)
		return true
	}
	return false
}
