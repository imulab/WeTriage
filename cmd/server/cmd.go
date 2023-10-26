package server

import (
	"absurdlab.io/WeTriage/cmd/internal"
	"absurdlab.io/WeTriage/internal/httpx"
	"absurdlab.io/WeTriage/internal/stringx"
	"absurdlab.io/WeTriage/route"
	"absurdlab.io/WeTriage/topic"
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Command() *cli.Command {
	conf := new(config)

	return &cli.Command{
		Name:   "server",
		Usage:  "Starts the callback server",
		Flags:  conf.flags(),
		Action: func(c *cli.Context) error { return runApp(c.Context, conf) },
	}
}

func runApp(ctx context.Context, conf *config) error {
	return fx.New(
		fx.NopLogger,
		fx.Supply(conf),
		fx.Provide(
			newLogger,
			newMqttClient,
			conf.toTopicProperties,
			conf.toRouteProperties,
			topic.NewTriageStrategies,
			route.NewEchoHandler,
			route.NewBusinessDataHandler,
		),
		fx.Invoke(
			func(
				logger *zerolog.Logger,
				echoHandler *route.EchoHandler,
				bizHandler *route.BusinessDataHandler,
			) error {

				mux := http.NewServeMux()
				{
					mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						httpx.WriteText(w, http.StatusOK, "WeTriage is running.")
					})

					mux.HandleFunc(conf.Path, func(w http.ResponseWriter, r *http.Request) {
						switch r.Method {
						case http.MethodGet:
							echoHandler.ServeHTTP(w, r)
						case http.MethodPost:
							bizHandler.ServeHTTP(w, r)
						default:
							httpx.WriteText(w, http.StatusMethodNotAllowed, "Method not allowed.")
						}
					})
				}

				log.Info().
					Int("port", conf.Port).
					Msg("Listening for incoming requests.")

				return http.ListenAndServe(
					fmt.Sprintf(":%d", conf.Port),
					httpx.NoPanic(logger)(mux),
				)
			},
		),
	).Start(ctx)
}

func newLogger(c *config) *zerolog.Logger {
	var lvl = zerolog.InfoLevel
	if c.Debug {
		lvl = zerolog.DebugLevel
	}

	logger := zerolog.New(os.Stderr).
		Level(lvl).
		With().
		Timestamp().
		Logger()

	return &logger
}

func newMqttClient(c *config, logger *zerolog.Logger) (*autopaho.ConnectionManager, error) {
	brokerUrl, err := url.Parse(c.MqttUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid mqtt broker url: %w", err)
	}

	pahoConnUp := make(chan struct{})

	var (
		errorLogger             = internal.NewPahoZeroLogger(logger)
		debugLogger paho.Logger = paho.NOOPLogger{}
	)
	if c.Debug {
		debugLogger = errorLogger
	}

	cm, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerUrl},
		KeepAlive:         60,
		ConnectRetryDelay: time.Millisecond,
		ConnectTimeout:    15 * time.Second,
		OnConnectionUp:    func(_ *autopaho.ConnectionManager, _ *paho.Connack) { pahoConnUp <- struct{}{} },
		Debug:             debugLogger,
		PahoDebug:         debugLogger,
		PahoErrors:        errorLogger,
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("WeTriageServer@%s", stringx.RandAlphaNumeric(6)),
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-pahoConnUp:
	case <-time.After(1 * time.Minute):
		return nil, errors.New("timeout exceeded when connecting to mqtt broker")
	}

	return cm, nil
}
