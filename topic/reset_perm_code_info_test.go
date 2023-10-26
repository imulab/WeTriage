package topic

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResetPermanentCodeInfoTriageStrategy(t *testing.T) {
	strategy := resetPermanentCodeInfoTriageStrategy{}

	const data = `
<xml>
	<SuiteId><![CDATA[dk4asffe9xxx4c0f4c]]></SuiteId>
	<AuthCode><![CDATA[AUTHCODE]]></AuthCode>
	<InfoType><![CDATA[reset_permanent_code]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
</xml>
`

	var features Features
	if err := xml.Unmarshal([]byte(data), &features); assert.NoError(t, err) {
		return
	}

	assert.True(t, strategy.Accepts(&features))

	topic, err := strategy.ParseXML([]byte(data))
	if assert.NoError(t, err) {
		assert.Equal(t, ResetPermanentCodeInfoName, topic.Name())
	}
}
