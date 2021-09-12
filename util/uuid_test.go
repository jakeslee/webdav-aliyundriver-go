package util

import "testing"

func TestNextIdStr(t *testing.T) {
	// 雪花测试
	println(NextIdStr())
	println(NextId())
}
