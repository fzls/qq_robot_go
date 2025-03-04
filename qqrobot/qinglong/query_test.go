package qinglong

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 2021/12/13 20:33 by fzls

func TestQueryCookieInfo(t *testing.T) {
	var info *JdCookieInfo

	// 无参数
	info = QueryCookieInfo("")
	assert.Nil(t, info)

	// pin
	info = QueryCookieInfo("pin_1")
	assert.NotNil(t, info)

	// 备注
	info = QueryCookieInfo("测试账号-1")
	assert.NotNil(t, info)

	// 仅env中存在的账号，使用pin
	info = QueryCookieInfo("pin_3")
	assert.NotNil(t, info)

	// 不存在的账号
	info = QueryCookieInfo("not exists")
	assert.Nil(t, info)
}

func TestQueryChartPath(t *testing.T) {
	info := QueryCookieInfo("pin_1")
	chartPath := QueryChartPath(info)
	expected, _ := filepath.Abs(getPath("log/.bean_chart/chart_pin_1.jpeg"))
	assert.Equal(t, expected, chartPath)
}

func TestQuerySummary(t *testing.T) {
	info := QueryCookieInfo("pin_1")
	assert.NotEmpty(t, QuerySummary(info))

	info = QueryCookieInfo("pin_3")
	assert.Empty(t, QuerySummary(info))

	// 即使所有日志都不包含农场，也应返回最后一个不为空的日志解析结果
	info = QueryCookieInfo(url.QueryEscape("pin_5"))
	assert.NotEmpty(t, QuerySummary(info))
}

func TestQueryCookieExpired(t *testing.T) {
	info := QueryCookieInfo("pin_1")
	assert.NotEmpty(t, QueryCookieExpired(info))

	info = QueryCookieInfo("pin_2")
	assert.NotEmpty(t, QueryCookieExpired(info))

	info = QueryCookieInfo("pin_3")
	assert.NotEmpty(t, QueryCookieExpired(info))

	info = QueryCookieInfo(url.QueryEscape("中文pin"))
	assert.NotEmpty(t, QueryCookieExpired(info))

	info = QueryCookieInfo("pin_99999")
	assert.Empty(t, QueryCookieExpired(info))
}
