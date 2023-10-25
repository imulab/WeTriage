package printer

import (
	"absurdlab.io/WeTriage/cmd"
	"absurdlab.io/WeTriage/internal/stringx"
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"net/url"
	"os"
	"time"
)

func Command() *cli.Command {
	conf := new(config)

	return &cli.Command{
		Name:   "printer",
		Usage:  "Simple topic consumer that just prints all messages to console",
		Flags:  conf.flags(),
		Action: func(c *cli.Context) error { return runApp(c.Context, conf) },
	}
}

func runApp(ctx context.Context, conf *config) error {
	return fx.New(
		fx.NopLogger,
		fx.Supply(conf),
		fx.Provide(newMqttClient, newLogger),
		fx.Invoke(
			func(logger *zerolog.Logger, client *autopaho.ConnectionManager) error {
				logger.Info().Msg("Waiting for WeTriage messages.")
				<-client.Done()
				return nil
			},
		),
	).Start(ctx)
}

func newMqttClient(c *config, logger *zerolog.Logger) (*autopaho.ConnectionManager, error) {
	brokerUrl, err := url.Parse(c.MqttUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid mqtt broker url: %w", err)
	}

	pahoConnUp := make(chan struct{})

	ctx := context.Background()

	var (
		errorLogger             = cmd.NewPahoZeroLogger(logger)
		debugLogger paho.Logger = paho.NOOPLogger{}
	)
	if c.Debug {
		debugLogger = errorLogger
	}

	cm, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerUrl},
		KeepAlive:         60,
		ConnectRetryDelay: time.Millisecond,
		ConnectTimeout:    15 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			if _, subErr := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: "T/WeTriage/+"},
				},
			}); subErr != nil {
				panic(subErr)
			}
			pahoConnUp <- struct{}{}
		},
		Debug:      debugLogger,
		PahoDebug:  debugLogger,
		PahoErrors: errorLogger,
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("WeTriagePrinter@%s", stringx.RandAlphaNumeric(6)),
			Router: paho.NewSingleHandlerRouter(func(pub *paho.Publish) {
				logger.Info().
					Str("topic", pub.Topic).
					Int("qos", int(pub.QoS)).
					RawJSON("payload", pub.Payload).
					Msg("Received message")
			}),
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
func newLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
	return &logger
}
