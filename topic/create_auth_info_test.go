package topic

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAuthInfoTriageStrategy(t *testing.T) {
	strategy := createAuthInfoTriageStrategy{}

	const data = `
<xml>
	<SuiteId><![CDATA[ww4asffe9xxx4c0f4c]]></SuiteId>
	<AuthCode><![CDATA[AUTHCODE]]></AuthCode>
	<InfoType><![CDATA[create_auth]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<State><![CDATA[123]]></State>
</xml>
`

	var features Features
	if err := xml.Unmarshal([]byte(data), &features); assert.NoError(t, err) {
		return
	}

	assert.True(t, strategy.Accepts(&features))

	topic, err := strategy.ParseXML([]byte(data))
	if assert.NoError(t, err) {
		assert.Equal(t, CreateAuthInfoName, topic.Name())
		assert.Equal(t, "ww4asffe9xxx4c0f4c", topic.(*CreateAuthInfo).SuiteId)
		assert.Equal(t, "AUTHCODE", topic.(*CreateAuthInfo).AuthCode)
		assert.Equal(t, "123", topic.(*CreateAuthInfo).State)
	}
}
