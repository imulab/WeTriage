package route

import (
	"absurdlab.io/WeTriage/internal/httpx"
	"absurdlab.io/WeTriage/topic"
	"net/http"
)

type respondStrategy interface {
	supports(topic topic.Topic) bool
	respond(w http.ResponseWriter, r *http.Request, topic topic.Topic) error
}

type fallbackRespondStrategy struct{}

func (s fallbackRespondStrategy) supports(_ topic.Topic) bool {
	return true
}

func (s fallbackRespondStrategy) respond(w http.ResponseWriter, _ *http.Request, _ topic.Topic) error {
	httpx.WriteText(w, http.StatusOK, "")
	return nil
}

type successTextRespondStrategy struct{}

func (_ successTextRespondStrategy) supports(t topic.Topic) bool {
	switch t.Name() {
	case topic.SuiteTicketInfoName: // add others
		return true
	default:
		return false
	}
}

func (_ successTextRespondStrategy) respond(w http.ResponseWriter, _ *http.Request, _ topic.Topic) error {
	httpx.WriteText(w, http.StatusOK, "success")
	return nil
}
