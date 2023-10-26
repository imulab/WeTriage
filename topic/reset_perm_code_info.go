package topic

import "encoding/xml"

const (
	ResetPermanentCodeInfoName = "reset_permanent_code_info"
)

var (
	_ Topic          = (*ResetPermanentCodeInfo)(nil)
	_ TriageStrategy = (*resetPermanentCodeInfoTriageStrategy)(nil)
)

// ResetPermanentCodeInfo 重置永久授权码通知. InfoType is always reset_permanent_code
//
// https://developer.work.weixin.qq.com/document/path/97175
type ResetPermanentCodeInfo struct {
	SuiteId   string `xml:"SuiteId" json:"suite_id"`
	AuthCode  string `xml:"AuthCode" json:"auth_code"`
	InfoType  string `xml:"InfoType" json:"info_type"`
	Timestamp int64  `xml:"TimeStamp" json:"timestamp"`
}

func (_ ResetPermanentCodeInfo) Name() string {
	return ResetPermanentCodeInfoName
}

type resetPermanentCodeInfoTriageStrategy struct{}

func (_ resetPermanentCodeInfoTriageStrategy) Accepts(f *Features) bool {
	return f.InfoType == "reset_permanent_code"
}

func (_ resetPermanentCodeInfoTriageStrategy) ParseXML(data []byte) (Topic, error) {
	var message ResetPermanentCodeInfo
	if err := xml.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
