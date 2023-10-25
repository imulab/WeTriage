package topic

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuiteTicketInfoTriageStrategy(t *testing.T) {
	strategy := suiteTicketInfoTriageStrategy{}

	const data = `
	<xml>
		<SuiteId><![CDATA[ww4asffe99e54c0fxxxx]]></SuiteId>
		<InfoType><![CDATA[suite_ticket]]></InfoType>
		<TimeStamp>1403610513</TimeStamp>
		<SuiteTicket><![CDATA[asdfasfdasdfasdf]]></SuiteTicket>
	</xml>
`

	var features Features
	if err := xml.Unmarshal([]byte(data), &features); assert.NoError(t, err) {
		return
	}

	assert.True(t, strategy.Accepts(&features))

	topic, err := strategy.ParseXML([]byte(data))
	if assert.NoError(t, err) {
		assert.Equal(t, SuiteTicketInfoName, topic.Name())
	}
}
