package route

import (
	"absurdlab.io/WeTriage/internal/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEchoHandler_ServeHTTP(t *testing.T) {
	h, err := NewEchoHandler(
		&Properties{
			Token: "QDG6eK",
			AesEncodingKey: func() []byte {
				k, err := crypto.Base64Std.Decode([]byte("jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"))
				require.NoError(t, err)
				return k
			}(),
		},
		&zerolog.Logger{},
	)
	require.NoError(t, err)

	q := url.Values{
		"msg_signature": []string{"5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"},
		"timestamp":     []string{"1409659589"},
		"nonce":         []string{"263014780"},
		"echostr":       []string{"P9nAzCzyDtyTWESHep1vC5X9xho/qYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp+4RPcs8TgAE7OaBO+FZXvnaqQ=="},
	}

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "1616140317555161061", rw.Body.String())
}
