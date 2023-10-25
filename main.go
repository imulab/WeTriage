package main

import (
	"absurdlab.io/WeTriage/buildinfo"
	"absurdlab.io/WeTriage/cmd/printer"
	"absurdlab.io/WeTriage/cmd/server"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.Info().
		Str("version", buildinfo.Version).
		Str("revision", buildinfo.Revision).
		Time("compiled_at", buildinfo.CompiledAtTime()).
		Msg("WeTriage binary")

	app := &cli.App{
		Name:      "WeTriage",
		Usage:     "WeTriage is a service to triage WeCom callbacks. Incoming XML messages are identified based on their traits and converted to JSON before handing off to a pluggable handler (i.e. message broker). Downstream services will have knowledge of the message type and can parse them with ease.",
		Version:   buildinfo.Version,
		Compiled:  buildinfo.CompiledAtTime(),
		Copyright: "MIT",
		Authors: []*cli.Author{
			{Name: "Weinan Qiu", Email: "davidiamyou@gmail.com"},
		},
		Commands: []*cli.Command{
			server.Command(),
			printer.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app.")
	}
}
