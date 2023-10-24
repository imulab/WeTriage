package main

import (
	"absurdlab.io/WeTriage/buildinfo"
	"absurdlab.io/WeTriage/internal/crypto"
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

type config struct {
	Port    int
	Debug   bool
	Path    string
	Token   string
	AesKey  string
	Topics  cli.StringSlice
	MqttUrl string
}

func (c *config) flags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Usage:       "Port to listen on",
			Value:       8080,
			EnvVars:     []string{"WT_PORT"},
			Destination: &c.Port,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug mode",
			EnvVars:     []string{"WT_DEBUG"},
			Destination: &c.Debug,
		},
		&cli.StringFlag{
			Name:        "path",
			Usage:       "Path to receive callbacks",
			Value:       "/callback",
			EnvVars:     []string{"WT_PATH"},
			Destination: &c.Path,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Token used to verify the authenticity of the callback",
			EnvVars:     []string{"WT_TOKEN"},
			Destination: &c.Token,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "aes-key",
			Usage:       "Base64 encoded AES key used to decrypt the callback",
			EnvVars:     []string{"WT_AES_KEY"},
			Destination: &c.AesKey,
			Required:    true,
		},
		&cli.StringSliceFlag{
			Name:        "topic",
			Usage:       "Enable handling for a supported topic",
			Destination: &c.Topics,
			Aliases:     []string{"t"},
		},
		&cli.StringFlag{
			Name:        "mqtt-url",
			Usage:       "MQTT broker url",
			EnvVars:     []string{"WT_MQTT_URL"},
			Destination: &c.MqttUrl,
			Required:    true,
		},
	}
}

func (c *config) newLogger() *zerolog.Logger {
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

func (c *config) newMqttClient() (*autopaho.ConnectionManager, error) {
	brokerUrl, err := url.Parse(c.MqttUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid mqtt broker url: %w", err)
	}

	timeout := 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pahoConnUp := make(chan struct{})

	cm, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerUrl},
		KeepAlive:         60,
		ConnectRetryDelay: time.Millisecond,
		ConnectTimeout:    15 * time.Second,
		OnConnectionUp:    func(_ *autopaho.ConnectionManager, _ *paho.Connack) { pahoConnUp <- struct{}{} },
		Debug:             paho.NOOPLogger{},
		PahoDebug:         paho.NOOPLogger{},
		PahoErrors:        paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("WeTriage@%s", stringx.RandAlphaNumeric(6)),
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-pahoConnUp:
	case <-time.After(timeout):
		return nil, errors.New("timeout exceeded when connecting to mqtt broker")
	}

	return cm, nil
}

func (c *config) toTopicProperties() *topic.Properties {
	return &topic.Properties{EnabledTopics: c.Topics.Value()}
}

func (c *config) toRouteProperties() (*route.Properties, error) {
	keyBytes, err := crypto.Base64Std.Decode([]byte(c.AesKey))
	if err != nil {
		return nil, fmt.Errorf("invalid aes encryption key encoding: %w", err)
	}

	return &route.Properties{
		Token:          c.Token,
		AesEncodingKey: keyBytes,
	}, nil
}

func runApp(ctx context.Context, conf *config) error {
	return fx.New(
		fx.NopLogger,
		fx.Provide(
			conf.newLogger,
			conf.newMqttClient,
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
				bizHandler route.BusinessDataHandler,
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

func main() {
	log.Info().
		Str("version", buildinfo.Version).
		Str("revision", buildinfo.Revision).
		Time("compiled_at", buildinfo.CompiledAtTime()).
		Msg("Starting WeTriage service.")

	conf := new(config)

	app := &cli.App{
		Name:      "WeTriage",
		Usage:     "WeTriage is a service to triage WeCom callbacks. Incoming XML messages are identified based on their traits and converted to JSON before handing off to a pluggable handler (i.e. message broker). Downstream services will have knowledge of the message type and can parse them with ease.",
		Version:   buildinfo.Version,
		Compiled:  buildinfo.CompiledAtTime(),
		Copyright: "MIT",
		Authors: []*cli.Author{
			{Name: "Weinan Qiu", Email: "davidiamyou@gmail.com"},
		},
		Flags:  conf.flags(),
		Action: func(cc *cli.Context) error { return runApp(cc.Context, conf) },
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app.")
	}
}
