package topic

import "encoding/xml"

const (
	CreateAuthInfoName = "create_auth_info"
)

var (
	_ Topic          = (*CreateAuthInfo)(nil)
	_ TriageStrategy = (*createAuthInfoTriageStrategy)(nil)
)

// CreateAuthInfo 授权成功通知. InfoType is always create_auth.
//
// https://developer.work.weixin.qq.com/document/path/97174
type CreateAuthInfo struct {
	SuiteId   string `xml:"SuiteId" json:"suite_id"`
	AuthCode  string `xml:"AuthCode" json:"auth_code"`
	InfoType  string `xml:"InfoType" json:"info_type"`
	Timestamp int64  `xml:"TimeStamp" json:"timestamp"`
	State     string `xml:"State" json:"state"`
}

func (_ CreateAuthInfo) Name() string {
	return CreateAuthInfoName
}

type createAuthInfoTriageStrategy struct{}

func (_ createAuthInfoTriageStrategy) Accepts(f *Features) bool {
	return f.InfoType == "create_auth"
}

func (_ createAuthInfoTriageStrategy) ParseXML(data []byte) (Topic, error) {
	var message CreateAuthInfo
	if err := xml.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
