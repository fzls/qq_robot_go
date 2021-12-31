package qinglong

import (
	"io/ioutil"
	"strings"

	"github.com/Mrs4s/go-cqhttp/global"
)

func parseCookieExpired(info *JdCookieInfo, logFilePath string) (result string, isLogComplete bool) {
	if !global.PathExists(logFilePath) {
		return "", false
	}

	contentBytes, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		return "", false
	}
	content := string(contentBytes)

	// 判断这个日志是否运行完整
	isLogComplete = strings.Contains(content, "开始发送通知...")

	prefixToRemove := " : "
	prefix := prefixToRemove + info.QueryUnescapedPtPin()
	suffix := "\n\n"

	// 定位前缀
	prefixIndex := strings.Index(content, prefix)
	if prefixIndex == -1 {
		return "", isLogComplete
	}
	prefixIndex += len(prefixToRemove)

	// 定位后缀
	suffixIndex := strings.Index(content[prefixIndex:], suffix)
	if suffixIndex == -1 {
		return "", isLogComplete
	}
	suffixIndex += prefixIndex

	result = content[prefixIndex:suffixIndex]

	return result, isLogComplete
}
