package topic

import "encoding/xml"

const (
	CancelAuthInfoName = "cancel_auth_info"
)

var (
	_ Topic          = (*CancelAuthInfo)(nil)
	_ TriageStrategy = (*cancelAuthInfoTriageStrategy)(nil)
)

// CancelAuthInfo 取消授权通知. InfoType is always cancel_auth.
//
// https://developer.work.weixin.qq.com/document/path/97174#%E5%8F%96%E6%B6%88%E6%8E%88%E6%9D%83%E9%80%9A%E7%9F%A5
type CancelAuthInfo struct {
	SuiteId    string `xml:"SuiteId" json:"suite_id"`
	InfoType   string `xml:"InfoType" json:"info_type"`
	Timestamp  int64  `xml:"TimeStamp" json:"timestamp"`
	AuthCorpId string `xml:"AuthCorpId" json:"auth_corp_id"`
}

func (_ CancelAuthInfo) Name() string {
	return CancelAuthInfoName
}

type cancelAuthInfoTriageStrategy struct{}

func (_ cancelAuthInfoTriageStrategy) Accepts(f *Features) bool {
	return f.InfoType == "cancel_auth"
}

func (_ cancelAuthInfoTriageStrategy) ParseXML(data []byte) (Topic, error) {
	var message CancelAuthInfo
	if err := xml.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
