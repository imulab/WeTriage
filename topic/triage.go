package topic

import (
	"fmt"
	"github.com/rs/zerolog"
)

// NewTriageStrategies returns a list of TriageStrategy implementations that are enabled by the given topic names.
func NewTriageStrategies(props *Properties) ([]TriageStrategy, error) {
	if len(props.EnabledTopics) == 0 {
		return nil, ErrNoTopicEnabled
	}

	var strategies []TriageStrategy
	for _, each := range props.EnabledTopics {
		switch each {
		case SuiteTicketInfoName:
			strategies = append(strategies, &suiteTicketInfoTriageStrategy{})
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupported, each)
		}
	}

	return strategies, nil
}

// TriageStrategy abstracts the logic of identifying a message type from its features.
type TriageStrategy interface {
	// Accepts returns true if the features shown by a message is recognized by this implementation.
	Accepts(f *Features) bool
	// ParseXML parses the given data into the message type supported by this implementation.
	ParseXML(data []byte) (Topic, error)
}

// Features collects identifying fields from a variety of callback message formats. By checking the fields known
// to a certain message format, its TriageStrategy can determine the type of the message and parse it.
type Features struct {
	InfoType     string `xml:"InfoType"`
	MsgType      string `xml:"MsgType"`
	Event        string `xml:"Event"`
	ChangeType   string `xml:"ChangeType"`
	BatchJobType string `xml:"BatchJob>JobType"`
	AuthType     string `xml:"AuthType"`
}

func (f *Features) MarshalZerologObject(e *zerolog.Event) {
	if len(f.InfoType) > 0 {
		e = e.Str("info_type", f.InfoType)
	}

	if len(f.MsgType) > 0 {
		e = e.Str("msg_type", f.MsgType)
	}

	if len(f.Event) > 0 {
		e = e.Str("event", f.Event)
	}

	if len(f.ChangeType) > 0 {
		e = e.Str("change_type", f.ChangeType)
	}

	if len(f.BatchJobType) > 0 {
		e = e.Str("batch_job_type", f.BatchJobType)
	}

	if len(f.AuthType) > 0 {
		e = e.Str("auth_type", f.AuthType)
	}
}
