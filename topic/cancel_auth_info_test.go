package topic

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCancelAuthInfoTriageStrategy(t *testing.T) {
	strategy := cancelAuthInfoTriageStrategy{}

	const data = `
<xml>
	<SuiteId><![CDATA[ww4asffe99e54cxxxx]]></SuiteId>
	<InfoType><![CDATA[cancel_auth]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<AuthCorpId><![CDATA[wxf8b4f85fxx794xxx]]></AuthCorpId>
</xml>
`

	var features Features
	if err := xml.Unmarshal([]byte(data), &features); assert.NoError(t, err) {
		return
	}

	assert.True(t, strategy.Accepts(&features))

	topic, err := strategy.ParseXML([]byte(data))
	if assert.NoError(t, err) {
		assert.Equal(t, CancelAuthInfoName, topic.Name())
		assert.Equal(t, "ww4asffe9xxx4c0f4c", topic.(*CancelAuthInfo).SuiteId)
		assert.Equal(t, "cancel_auth", topic.(*CancelAuthInfo).InfoType)
		assert.Equal(t, int64(1403610513), topic.(*CancelAuthInfo).Timestamp)
		assert.Equal(t, "wxf8b4f85f3a794xxx", topic.(*CancelAuthInfo).AuthCorpId)
	}
}
