package route

import (
	"absurdlab.io/WeTriage/internal/crypto"
	"absurdlab.io/WeTriage/internal/httpx"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"sort"
	"strings"
)

func NewEchoHandler(props *Properties, logger *zerolog.Logger) (*EchoHandler, error) {
	return &EchoHandler{
		props:    props,
		aesCrypt: crypto.NewAesCbcPkcs7Padding(props.AesEncodingKey),
		logger:   logger,
	}, nil
}

type EchoHandler struct {
	props    *Properties
	aesCrypt *crypto.AesCbcCrypt
	logger   *zerolog.Logger
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.doServerHTTP(w, r); err != nil {
		httpx.WriteText(w, http.StatusBadRequest, err.Error())
		h.logger.Err(err).Msg("Echo respond error")
		return
	}

	h.logger.Info().Msg("Responded to echo")
}

func (h *EchoHandler) doServerHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		signature     = r.URL.Query().Get("msg_signature")
		encryptedEcho = r.URL.Query().Get("echostr")

		parts = []string{
			h.props.Token,
			r.URL.Query().Get("timestamp"),
			r.URL.Query().Get("nonce"),
			encryptedEcho,
		}
	)

	sort.Strings(parts)

	var sigBytes []byte
	if sigBytes, err = crypto.Encode(
		[]byte(strings.Join(parts, "")),
		crypto.SHA1,
		crypto.Hex,
	); err != nil {
		return
	}

	if string(sigBytes) != signature {
		return errors.New("signature mismatch")
	}

	var echoBytes []byte
	echoBytes, _, err = decodeEncryptedMessage(h.aesCrypt, encryptedEcho)
	if err != nil {
		return
	}

	httpx.WriteText(w, http.StatusOK, string(echoBytes))

	return
}
