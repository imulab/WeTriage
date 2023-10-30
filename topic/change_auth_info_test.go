package topic

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeAuthInfoTriageStrategy(t *testing.T) {
	strategy := changeAuthInfoTriageStrategy{}

	const data = `
<xml>
	<SuiteId><![CDATA[ww4asffe99exxx0f4c]]></SuiteId>
	<InfoType><![CDATA[change_auth]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<State><![CDATA[abc]]></State>
</xml>
`

	var features Features
	if err := xml.Unmarshal([]byte(data), &features); assert.NoError(t, err) {
		return
	}

	assert.True(t, strategy.Accepts(&features))

	topic, err := strategy.ParseXML([]byte(data))
	if assert.NoError(t, err) {
		assert.Equal(t, ChangeAuthInfoName, topic.Name())
		assert.Equal(t, "ww4asffe9xxx4c0f4c", topic.(*ChangeAuthInfo).SuiteId)
		assert.Equal(t, "123", topic.(*ChangeAuthInfo).State)
		assert.Equal(t, "change_auth", topic.(*ChangeAuthInfo).InfoType)
		assert.Equal(t, int64(1403610513), topic.(*ChangeAuthInfo).Timestamp)
		assert.Equal(t, "wxf8b4f85f3a794e77", topic.(*ChangeAuthInfo).AuthCorpId)
	}
}
