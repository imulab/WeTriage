package server

import (
	"absurdlab.io/WeTriage/internal/crypto"
	"absurdlab.io/WeTriage/route"
	"absurdlab.io/WeTriage/topic"
	"fmt"
	"github.com/urfave/cli/v2"
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
