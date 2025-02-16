package qinglong

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJdCookie(t *testing.T) {
	ptPinToCookieInfo, err := ParseJdCookie()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(ptPinToCookieInfo))

	// 普通账号
	ck1, ok := ptPinToCookieInfo["pin_1"]
	assert.True(t, ok)
	assert.Equal(t, "测试账号-1", ck1.Remark)

	// remarks中直接写备注，不带 remark=前缀 和 ;后缀
	ck2, ok := ptPinToCookieInfo["pin_2"]
	assert.True(t, ok)
	assert.Equal(t, "测试账号-2", ck2.Remark)

	// 仅在env.sh中存在的账号
	ck3, ok := ptPinToCookieInfo["pin_3"]
	assert.True(t, ok)
	assert.Equal(t, "pin_3", ck3.Remark)

	// 不存在的账号
	_, ok = ptPinToCookieInfo["pin_4"]
	assert.False(t, ok)
}

func Test_parseEnvDB(t *testing.T) {
	ptPinToCookieInfo, err := parseEnvDB()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ptPinToCookieInfo))

	assert.NotNil(t, ptPinToCookieInfo["pin_1"])
	assert.Nil(t, ptPinToCookieInfo["pin_3"])
}

func Test_parseEnvSh(t *testing.T) {
	ptPinToIndex, err := parseEnvSh()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(ptPinToIndex))

	assert.Equal(t, 1, ptPinToIndex["pin_1"])
	assert.Equal(t, 3, ptPinToIndex["pin_3"])

	_, ok := ptPinToIndex["pin_4"]
	assert.False(t, ok)
}

func Test_getPtPin(t *testing.T) {
	assert.Equal(t, "YYYY", getPtPin("pt_key=XXXX;pt_pin=YYYY;"))
	assert.Equal(t, "", getPtPin(""))
	assert.Equal(t, "YYYY=1", getPtPin("pt_key=XXXX;pt_pin=YYYY=1;"))
	assert.Equal(t, "", getPtPin("pt_key=XXXX;pt_pinYYYY;"))
}

func Test_getCookie(t *testing.T) {
	assert.Equal(t, "YYYY", getCookie("pt_key=XXXX;pt_pin=YYYY;", "pt_pin"))
}

func Test_getRemark(t *testing.T) {
	assert.Equal(t, "test", getRemark("remark=test;"))
	assert.Equal(t, "test", getRemark("remark=test;a=b"))
	assert.Equal(t, "test", getRemark("remark=test;a=b;"))
	assert.Equal(t, "remark=test", getRemark("remark=remark=test;"))
	assert.Equal(t, "test", getRemark("test"))
	assert.Equal(t, "", getRemark(""))
	assert.Equal(t, "test", getRemark("remark=test;@@UID_xxxxxx"))
	assert.Equal(t, "test", getRemark("test@@UID_xxxxxx"))
	assert.Equal(t, "test", getRemark("test@@1640780099690@@UID_xxxxxx"))
}

func TestEnvDBEntry_CreateTime(t *testing.T) {
	expectedCreateTimeStamp := int64(1636536508)
	withoutExtraRemark := &EnvDBEntry{
		Created:   expectedCreateTimeStamp * 1000,
		Timestamp: "Wed Nov 10 2021 17:28:28 GMT+0800 (中国标准时间)",
		Remarks:   "test",
	}
	assert.Equal(t, expectedCreateTimeStamp, withoutExtraRemark.CreateTime())
	assert.Equal(t, expectedCreateTimeStamp, withoutExtraRemark.UpdateTime())

	expectedNvjdcUpdateTimeStamp := int64(1636536538)
	withExtraRemark := &EnvDBEntry{
		Created:   1636536508000,
		Timestamp: "Wed Nov 10 2021 17:28:28 GMT+0800 (中国标准时间)",
		Remarks:   fmt.Sprintf("test@@%v@@UID_xxxxxx", expectedNvjdcUpdateTimeStamp*1000),
	}
	assert.Equal(t, expectedCreateTimeStamp, withExtraRemark.CreateTime())
	assert.Equal(t, expectedNvjdcUpdateTimeStamp, withExtraRemark.UpdateTime())
}
