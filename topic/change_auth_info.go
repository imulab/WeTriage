package topic

import "encoding/xml"

const (
	ChangeAuthInfoName = "change_auth_info"
)

var (
	_ Topic          = (*ChangeAuthInfo)(nil)
	_ TriageStrategy = (*changeAuthInfoTriageStrategy)(nil)
)

// ChangeAuthInfo 授权变更通知. InfoType is always change_auth.
//
// https://developer.work.weixin.qq.com/document/path/97174#%E5%8F%98%E6%9B%B4%E6%8E%88%E6%9D%83%E9%80%9A%E7%9F%A5
type ChangeAuthInfo struct {
	SuiteId    string `xml:"SuiteId" json:"suite_id"`
	InfoType   string `xml:"InfoType" json:"info_type"`
	Timestamp  int64  `xml:"TimeStamp" json:"timestamp"`
	AuthCorpId string `xml:"AuthCorpId" json:"auth_corp_id"`
	State      string `xml:"State" json:"state"`
}

func (_ ChangeAuthInfo) Name() string {
	return ChangeAuthInfoName
}

type changeAuthInfoTriageStrategy struct{}

func (_ changeAuthInfoTriageStrategy) Accepts(f *Features) bool {
	return f.InfoType == "change_auth"
}

func (_ changeAuthInfoTriageStrategy) ParseXML(data []byte) (Topic, error) {
	var message ChangeAuthInfo
	if err := xml.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
