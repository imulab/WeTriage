package topic

import "encoding/xml"

const (
	SuiteTicketInfoName = "suite_ticket_info"
)

var (
	_ Topic          = (*SuiteTicketInfo)(nil)
	_ TriageStrategy = (*suiteTicketInfoTriageStrategy)(nil)
)

// SuiteTicketInfo 推送suite_ticket. InfoType is always suite_ticket.
//
// https://developer.work.weixin.qq.com/document/path/97173
type SuiteTicketInfo struct {
	SuiteId     string `xml:"SuiteId" json:"suite_id"`
	InfoType    string `xml:"InfoType" json:"info_type"`
	Timestamp   int64  `xml:"TimeStamp" json:"timestamp"`
	SuiteTicket string `xml:"SuiteTicket" json:"suite_ticket"`
}

func (_ SuiteTicketInfo) Name() string {
	return SuiteTicketInfoName
}

type suiteTicketInfoTriageStrategy struct{}

func (_ suiteTicketInfoTriageStrategy) Accepts(f *Features) bool {
	return f.InfoType == "suite_ticket"
}

func (_ suiteTicketInfoTriageStrategy) ParseXML(data []byte) (Topic, error) {
	var message SuiteTicketInfo
	if err := xml.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
