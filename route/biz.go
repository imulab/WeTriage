package route

import (
	"absurdlab.io/WeTriage/internal/crypto"
	"absurdlab.io/WeTriage/internal/httpx"
	"absurdlab.io/WeTriage/topic"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"net/http"
	"sort"
	"strings"
	"time"
)

func NewBusinessDataHandler(
	props *Properties,
	triageStrategies []topic.TriageStrategy,
	logger *zerolog.Logger,
	mqtt *autopaho.ConnectionManager,
) *BusinessDataHandler {
	return &BusinessDataHandler{
		props:            props,
		aesCrypt:         crypto.NewAesCbcPkcs7Padding(props.AesEncodingKey),
		mqtt:             mqtt,
		triageStrategies: triageStrategies,
		logger:           logger,
		respondStrategies: []respondStrategy{
			successTextRespondStrategy{},
			fallbackRespondStrategy{},
		},
	}
}

type BusinessDataHandler struct {
	props             *Properties
	aesCrypt          *crypto.AesCbcCrypt
	mqtt              *autopaho.ConnectionManager
	triageStrategies  []topic.TriageStrategy
	respondStrategies []respondStrategy
	logger            *zerolog.Logger
}

func (h *BusinessDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.doServeHTTP(w, r); err != nil {
		httpx.WriteText(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *BusinessDataHandler) doServeHTTP(w http.ResponseWriter, r *http.Request) error {
	messageBytes, err := h.verifyAndDecodeMessageBytes(r)
	if err != nil {
		return err
	}

	message, err := h.triageMessage(messageBytes)
	if err != nil {
		return err
	}

	if err = h.publishMessage(r.Context(), message); err != nil {
		return err
	}

	if err = h.respond(w, r, message); err != nil {
		return err
	}

	h.logger.Info().Str("topic", message.Name()).Msg("Callback processed")

	return nil
}

func (h *BusinessDataHandler) verifyAndDecodeMessageBytes(r *http.Request) ([]byte, error) {
	var body = struct {
		ToUserName string `xml:"ToUserName"`
		AgentID    string `xml:"AgentId"`
		Encrypt    string `xml:"Encrypt"`
	}{}
	if err := xml.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	parts := []string{h.props.Token, r.URL.Query().Get("timestamp"), r.URL.Query().Get("nonce"), body.Encrypt}
	sort.Strings(parts)

	if sigBytes, encErr := crypto.Encode(
		[]byte(strings.Join(parts, "")),
		crypto.SHA1,
		crypto.Hex,
	); encErr != nil {
		return nil, encErr
	} else if string(sigBytes) != r.URL.Query().Get("msg_signature") {
		return nil, errors.New("signature mismatch")
	}

	if messageBytes, _, err := decodeEncryptedMessage(h.aesCrypt, body.Encrypt); err != nil {
		return nil, err
	} else {
		return messageBytes, nil
	}
}

func (h *BusinessDataHandler) triageMessage(messageBytes []byte) (topic.Topic, error) {
	var features topic.Features
	if err := xml.Unmarshal(messageBytes, &features); err != nil {
		return nil, err
	}

	triageStrategy, hasTriageStrategy := lo.Find(h.triageStrategies, func(s topic.TriageStrategy) bool { return s.Accepts(&features) })
	if !hasTriageStrategy {
		h.logger.Debug().
			Object("features", &features).
			Msg("No topic.TriageStrategy found for message")
		return nil, topic.ErrUnsupported
	}

	if message, err := triageStrategy.ParseXML(messageBytes); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}

func (h *BusinessDataHandler) publishMessage(ctx context.Context, message topic.Topic) error {
	var payload = struct {
		Id        string      `json:"id"`
		CreatedAt int64       `json:"created_at"`
		Topic     string      `json:"topic"`
		Content   interface{} `json:"content"`
	}{
		Id:        uuid.NewString(),
		CreatedAt: time.Now().Unix(),
		Topic:     message.Name(),
		Content:   message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err = h.mqtt.Publish(ctx, &paho.Publish{
		QoS:        2,
		Topic:      fmt.Sprintf("T/WeTriage/%s", message.Name()),
		Properties: &paho.PublishProperties{ContentType: "application/json"},
		Payload:    payloadBytes,
	}); err != nil {
		h.logger.Debug().Err(err).Msg("Failed to publish message to mqtt broker")
		return err
	}

	return nil
}

func (h *BusinessDataHandler) respond(w http.ResponseWriter, r *http.Request, message topic.Topic) error {
	responder, _ := lo.Find(h.respondStrategies, func(item respondStrategy) bool {
		return item.supports(message)
	})

	return responder.respond(w, r, message)
}
